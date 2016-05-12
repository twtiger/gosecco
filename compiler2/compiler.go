package compiler2

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/twtiger/gosecco/constants"
	"github.com/twtiger/gosecco/tree"

	"golang.org/x/sys/unix"
)

// This compiler runs in three stages
// - generate base code
// - do peephole optimization
// - resolve all labels and jumps

// TODO: handle boolean literal at top level
// TODO: Fixup peephole optimization
// TODO: Fixup jumps
// TODO: handle full compile of rules, not just the expression
// TODO: put together more than one rule
// TODO: add the prefix and postfix
// TODO: fix all potential errors

// Compile will take a parsed policy and generate an optimized sock filter for that policy
// The policy is assumed to have been unified and simplified before compilation starts -
// no unresolved variables or calls should exist in the policy.
func Compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c := createCompilerContext()
	return c.compile(policy.Rules)
}

type label string

var positive = label("positive")
var negative = label("negative")
var noLabel = label("noLabel")

type compilerContext struct {
	result          []unix.SockFilter
	currentlyLoaded int
	stackTop        uint32
	// TODO we will need an unconditional jumps map as well.
	// these need to be fixed up if there is an unconditional jump inserted after.
	// these need to be fixed up if there is an unconditional jump inserted after.
	jts          map[label][]int
	jfs          map[label][]int
	uconds       map[label][]int
	labels       map[label]int
	labelCounter int
	// this will always be 0xFF in production, but it is injectable for testing.
	maxJumpSize                  int
	currentlyCompilingSyscall    string
	currentlyCompilingExpression tree.Expression
}

func createCompilerContext() *compilerContext {
	return &compilerContext{
		jts:             make(map[label][]int),
		jfs:             make(map[label][]int),
		uconds:          make(map[label][]int),
		labels:          make(map[label]int),
		maxJumpSize:     255,
		currentlyLoaded: -1,
	}
}

func (c *compilerContext) compile(rules []tree.Rule) ([]unix.SockFilter, error) {
	for _, r := range rules {
		c.compileRule(r)
	}

	// TODO at end of rules we should have a jump to the default action

	c.negativeAction()
	c.positiveAction()

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

func (c *compilerContext) checkCorrectSyscall(name string, setPosFlags bool, next label) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	if setPosFlags {
		c.opWithJumps(OP_JEQ_K, sys, positive, next)
	} else {
		c.opWithJumps(OP_JEQ_K, sys, noLabel, next)
	}
}

func (c *compilerContext) compileRule(r tree.Rule) {
	next := c.newLabel()
	neg := c.newLabel()

	_, isBoolLit := r.Body.(tree.BooleanLiteral)
	c.checkCorrectSyscall(r.Name, isBoolLit, next) // set JT flag to final ret_allow only if the rule is a boolean literal

	// These are useful for debugging and helpful error messages
	c.currentlyCompilingSyscall = r.Name
	c.currentlyCompilingExpression = r.Body

	c.compileExpression(r.Body, neg)

	c.labelHere(next)
	c.labelHere(neg)
}

func (c *compilerContext) positiveAction() {
	c.labelHere(positive)
	c.op(OP_RET_K, SECCOMP_RET_ALLOW)
}

func (c *compilerContext) negativeAction() {
	c.labelHere(negative)
	c.op(OP_RET_K, SECCOMP_RET_KILL)
}

func (c *compilerContext) op(code uint16, k uint32) {
	c.result = append(c.result, unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
}

func (c *compilerContext) compileExpression(x tree.Expression, neg label) {
	// Returns error
	isTopLevel := true
	compileBoolean(c, x, isTopLevel, positive, neg)
}

func (c *compilerContext) newLabel() label {
	result := fmt.Sprintf("generatedLabel%03d", c.labelCounter)
	c.labelCounter++
	return label(result)
}

func (c *compilerContext) registerJumps(index int, jt, jf label) {
	c.jts[jt] = append(c.jts[jt], index)
	c.jfs[jf] = append(c.jfs[jf], index)
}

func (c *compilerContext) fixMaxJumps(l label, elems []int, isPos bool) {
	to := len(c.result)
	for _, from := range elems {
		if (to-from)-1 > c.maxJumpSize {
			c.longJump(from, isPos, l)
		}
	}
}

func (c *compilerContext) labelHere(l label) {
	//fmt.Println("jump trues", c.jts)
	//fmt.Println("jump false", c.jfs)
	//fmt.Println("unconditional jumps", c.uconds)
	//fmt.Println("labels", c.labels)
	//fmt.Println("+++++++++++++++++++++++++++++++++++")

	jts, jfs := c.jts[l], c.jfs[l]

	c.fixMaxJumps(l, jts, true)
	c.fixMaxJumps(l, jfs, false)
	c.labels[l] = len(c.result)
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
