package compiler2

import (
	"errors"
	"fmt"
	"sort"
	"syscall"

	"github.com/twtiger/gosecco/constants"
	"github.com/twtiger/gosecco/tree"

	"golang.org/x/sys/unix"
)

// TODO: add the prefix and postfix
// TODO: compare go-seccomp and gosecco policy evaluation

// Compile will take a parsed policy and generate an optimized sock filter for that policy
// The policy is assumed to have been unified and simplified before compilation starts -
// no unresolved variables or calls should exist in the policy.
func Compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c := createCompilerContext()
	return c.compile(policy)
}

type label string

type compilerContext struct {
	result                       []unix.SockFilter
	currentlyLoaded              int
	stackTop                     uint32
	jts                          *jumpMap
	jfs                          *jumpMap
	uconds                       *jumpMap
	labels                       map[label]int
	labelCounter                 int
	defaultPositive              string
	defaultNegative              string
	actions                      map[string]label
	maxJumpSize                  int // this will always be 0xFF in production, but can be injected for testing.
	currentlyCompilingSyscall    string
	currentlyCompilingExpression tree.Expression
}

func createCompilerContext() *compilerContext {
	return &compilerContext{
		jts:             createJumpMap(),
		jfs:             createJumpMap(),
		uconds:          createJumpMap(),
		labels:          make(map[label]int),
		actions:         make(map[string]label),
		maxJumpSize:     255,
		currentlyLoaded: -1,
	}
}

func (c *compilerContext) setDefaults(positive, negative string) {
	if positive != "" {
		c.defaultPositive = positive
	} else {
		c.defaultPositive = defaultPositive
	}

	if negative != "" {
		c.defaultNegative = negative
	} else {
		c.defaultNegative = defaultNegative
	}
}

func (c *compilerContext) getOrCreateAction(action string) label {
	l, lExists := c.actions[action]

	if lExists {
		return l
	} else {
		actionLabel := c.newLabel()
		c.actions[action] = actionLabel
		return actionLabel
	}
}

func sortActions(s map[string]label) []string {
	actionOrder := []string{}

	for k := range s {
		actionOrder = append(actionOrder, k)
	}

	sort.Strings(actionOrder)
	return actionOrder
}

func (c *compilerContext) compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c.setDefaults(policy.DefaultPositiveAction, policy.DefaultNegativeAction)

	for _, r := range policy.Rules {
		c.compileRule(r)
	}

	var defAction string

	if policy.DefaultPolicyAction == "" {
		defAction = defaultNegative
	} else {
		defAction = policy.DefaultPolicyAction
	}

	l := c.getOrCreateAction(defAction)
	c.unconditionalJumpTo(l) // Default action if we don't set this in the policy

	actionOrder := sortActions(c.actions)

	for _, k := range actionOrder {
		c.labelHere(c.actions[k])
		r := actionInstructions[k]
		c.op(OP_RET_K, r)
	}

	c.fixupJumps()

	return c.result, nil
}

func (c *compilerContext) loadAt(pos uint32) {
	if c.currentlyLoaded != int(pos) {
		c.op(OP_LOAD, pos)
		c.currentlyLoaded = int(pos)
	}
}

func (c *compilerContext) loadLiteral(lit uint32) {
	c.op(OP_LOAD_VAL, lit)
	c.currentlyLoaded = -1
}

const syscallNameIndex = 0

func (c *compilerContext) loadCurrentSyscall() {
	c.loadAt(syscallNameIndex)
}

func (c *compilerContext) checkCorrectSyscall(name string, next label) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	goesNowhere := c.newLabel()
	c.opWithJumps(OP_JEQ_K, sys, goesNowhere, next)
	c.labelHere(goesNowhere)
}

func (c *compilerContext) compileRule(r tree.Rule) {
	next := c.newLabel()

	pos, neg := c.compileActions(r.PositiveAction, r.NegativeAction)

	c.checkCorrectSyscall(r.Name, next) // set JT flag to final ret_allow only if the rule is a boolean literal

	// These are useful for debugging and helpful error messages
	c.currentlyCompilingSyscall = r.Name
	c.currentlyCompilingExpression = r.Body

	c.compileExpression(r.Body, pos, neg)

	c.labelHere(next)
}

func (c *compilerContext) compileActions(positiveAction string, negativeAction string) (label, label) {
	if positiveAction == "" {
		positiveAction = c.defaultPositive
	}

	if negativeAction == "" {
		negativeAction = c.defaultNegative
	}

	posActionLabel := c.getOrCreateAction(positiveAction)
	negActionLabel := c.getOrCreateAction(negativeAction)

	return posActionLabel, negActionLabel
}

func (c *compilerContext) op(code uint16, k uint32) {
	c.result = append(c.result, unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
}

func (c *compilerContext) compileExpression(x tree.Expression, pos, neg label) {
	// Returns error
	isTopLevel := true
	compileBoolean(c, x, isTopLevel, pos, neg)
}

func (c *compilerContext) newLabel() label {
	result := fmt.Sprintf("generatedLabel%03d", c.labelCounter)
	c.labelCounter++
	return label(result)
}

func (c *compilerContext) registerJumps(index int, jt, jf label) {
	c.jts.registerJump(jt, index)
	c.jfs.registerJump(jf, index)
}

func (c *compilerContext) fixMaxJumps(l label, j *jumpMap, isPos bool) {
	to := len(c.result)
	for _, from := range j.allJumpsTo(l) {
		if (to-from)-1 > c.maxJumpSize {
			c.longJump(from, isPos, l)
		}
	}
}

func (c *compilerContext) labelHere(l label) {
	c.fixMaxJumps(l, c.jts, true)
	c.fixMaxJumps(l, c.jfs, false)
	c.labels[l] = len(c.result)
}

func (c *compilerContext) unconditionalJumpTo(to label) {
	index := len(c.result)
	c.result = append(c.result, unix.SockFilter{
		Code: OP_JMP_K,
		Jt:   0,
		Jf:   0,
		K:    0,
	})
	c.uconds.registerJump(to, index)
}

func (c *compilerContext) opWithJumps(code uint16, k uint32, jt, jf label) {
	index := len(c.result)
	c.registerJumps(index, jt, jf)
	c.result = append(c.result, unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
}

func (c *compilerContext) jumpOnEq(val uint32, jt, jf label) {
	c.opWithJumps(OP_JEQ_K, val, jt, jf)
}

func (c *compilerContext) pushAToStack() error {
	if c.stackTop >= syscall.BPF_MEMWORDS {
		return errors.New("the expression is too complicated to compile. Please refer to the language documentation")
	}

	c.op(OP_STORE, c.stackTop)
	c.stackTop++
	return nil
}

func (c *compilerContext) popStackToX() error {
	if c.stackTop == 0 {
		return errors.New("popping from empty stack - this is likely a programmer error")
	}
	c.stackTop--
	c.op(OP_LOAD_MEM_X, c.stackTop)
	return nil
}
