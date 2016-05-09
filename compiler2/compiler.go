package compiler2

import (
	"errors"
	"fmt"

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
var noLabel = label("")

type compilerContext struct {
	result                       []unix.SockFilter
	currentlyLoaded              int
	stackTop                     uint32
	jts                          map[label][]int
	jfs                          map[label][]int
	labels                       map[label]int
	labelCounter                 int
	maxJumpSize                  int
	currentlyCompilingSyscall    string
	currentlyCompilingExpression tree.Expression
}

func createCompilerContext() *compilerContext {
	return &compilerContext{
		jts:             make(map[label][]int),
		jfs:             make(map[label][]int),
		labels:          make(map[label]int),
		maxJumpSize:     1,
		currentlyLoaded: -1,
	}
}

func (c *compilerContext) compile(rules []tree.Rule) ([]unix.SockFilter, error) {
	for _, r := range rules {
		c.compileRule(r)
	}

	// at end of rules we should have a jump to the default action

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
	_, isBoolLit := r.Body.(tree.BooleanLiteral)
	c.checkCorrectSyscall(r.Name, isBoolLit, next) // set JT flag to final ret_allow only if the rule is a boolean literal

	// These are useful for debugging and helpful error messages
	c.currentlyCompilingSyscall = r.Name
	c.currentlyCompilingExpression = r.Body

	c.compileExpression(r.Body)

	c.labelHere(next)
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

func (c *compilerContext) compileExpression(x tree.Expression) {
	// Returns error
	compileBoolean(c, x, true, positive, negative)
	c.fixupJumps()
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

func (c *compilerContext) labelHere(l label) {
	at := len(c.result)

	jts, jfs := c.jts[l], c.jfs[l]

	resultLabel := make([]int, 0, len(jts))
	for _, pos := range jts {
		if (at-pos)-1 > c.maxJumpSize {
			// insert a new JUMP, pointing at this label
			// We need a new thing that can fix up for direct jumps
			// everything after needs to be fixed up
			// We need to check that both jt and jf point correctly afterwards
			// specifically, jt should point to 0. But if jf is 0, it needs to be 1, etc.
			// maybe the generic algorithm can take care of this

			fmt.Printf("labelHere: %s at: %d\n", string(l), at)
			fmt.Printf("Blarg: %d\n", pos)
		}
		resultLabel = append(resultLabel, pos)
	}
	//	fmt.Printf("One: %#v  Two: %#v\n", jts, resultLabel)
	if jts != nil {
		jts = resultLabel
	}

	for _, pos := range jfs {
		if (at-pos)-1 >= c.maxJumpSize {
			//			fmt.Println("Blarg")
		}
	}

	// Check if this is a long jump
	// Then immediately insert the long jump, and remove the jump points
	// and then fix up everything.
	// We have to fixup all the other jump points directly afterwards

	c.labels[l] = at
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
	if c.stackTop >= 4294967295 {
		return errors.New("Stack limit reached")
	} else {
		c.op(OP_STORE, c.stackTop)
		c.stackTop++
		return nil
	}
}

func (c *compilerContext) popStackToX() error {
	c.stackTop--
	c.op(OP_LOAD_MEM_X, c.stackTop)
	return nil
}
