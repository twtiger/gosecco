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

// TODO: handle full compile of rules, not just the expression
// TODO: add the prefix and postfix
// TODO: fix all potential errors (no panics, we should check for errors)
// TODO: compare go-seccomp and gosecco policy evaluation

// Compile will take a parsed policy and generate an optimized sock filter for that policy
// The policy is assumed to have been unified and simplified before compilation starts -
// no unresolved variables or calls should exist in the policy.
func Compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c := createCompilerContext()
	return c.compile(policy.Rules)
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
	actions                      map[string]label
	maxJumpSize                  int // this is always be 0xFF in production, but can be injected for testing.
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

func (c *compilerContext) compile(rules []tree.Rule) ([]unix.SockFilter, error) {
	for _, r := range rules {
		c.compileRule(r)
	}

	// TODO: use default policy here instead of kill
	c.unconditionalJumpTo(c.actions["kill"])

	actionOrder := []string{}
	for k := range c.actions {
		actionOrder = append(actionOrder, k)
	}
	sort.Strings(actionOrder)

	for _, k := range actionOrder {
		c.labelHere(c.actions[k])
		switch k {
		case "allow":
			c.op(OP_RET_K, SECCOMP_RET_ALLOW)
		case "kill":
			c.op(OP_RET_K, SECCOMP_RET_KILL)
		case "trace":
			c.op(OP_RET_K, SECCOMP_RET_TRACE)
		}
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
		positiveAction = "allow"
	}

	if negativeAction == "" {
		negativeAction = "kill"
	}

	posActionLabel, positiveActionExists := c.actions[positiveAction]
	negActionLabel, negativeActionExists := c.actions[negativeAction]

	if !positiveActionExists {
		posActionLabel = c.newLabel()
		c.actions[positiveAction] = posActionLabel
	}

	if !negativeActionExists {
		negActionLabel = c.newLabel()
		c.actions[negativeAction] = negActionLabel
	}

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
