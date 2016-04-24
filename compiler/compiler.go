package compiler

import (
	"github.com/twtiger/gosecco/constants"
	"github.com/twtiger/gosecco/tree"
	"golang.org/x/sys/unix"
)

func newCompiler() *compiler {
	return &compiler{
		currentlyLoaded: -1,
		positiveLabels:  make(map[label][]labelInfo),
		negativeLabels:  make(map[label][]labelInfo),
	}
}

// Compile will take a parsed policy and generate an optimized sock filter for that policy
// The policy is assumed to have been unified and simplified before compilation starts -
// no unresolved variables or calls should exist in the policy.
func Compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c := newCompiler()
	c.compile(policy.Rules)
	return c.result, nil
}

type compiler struct {
	result          []unix.SockFilter
	currentlyLoaded int
	positiveLabels  map[label][]labelInfo
	negativeLabels  map[label][]labelInfo
}

func (c *compiler) compile(rules []tree.Rule) {
	for _, r := range rules {
		c.compileRule(r)
	}
	c.positiveAction(noLabel)
	c.negativeAction(noLabel)
}

func (c *compiler) compileExpression(x tree.Expression) {
	cv := &compilerVisitor{c, true, true, false}
	x.Accept(cv)
}

func (c *compiler) compileRule(r tree.Rule) {
	c.labelHere(negative)
	_, isBoolLit := r.Body.(tree.BooleanLiteral)
	c.checkCorrectSyscall(r.Name, isBoolLit) // set JT flag to final ret_allow only if the rule is a boolean literal
	c.compileExpression(r.Body)
}

const syscallNameIndex = 0

type kexInstruction struct {
	k uint16
	x uint16
}

var comparisonOps = map[tree.ComparisonType]kexInstruction{
	tree.EQL:  kexInstruction{k: JEQ_K, x: JEQ_X},
	tree.NEQL: kexInstruction{k: JEQ_K, x: JEQ_X},
	tree.GT:   kexInstruction{k: JEG_K, x: JEG_X},
	tree.GTE:  kexInstruction{k: JEGE_K, x: JEGE_X},
	tree.LT:   kexInstruction{k: JEG_K, x: JEG_X},
	tree.LTE:  kexInstruction{k: JEGE_K, x: JEGE_X},
}

var posVals = map[tree.ComparisonType]bool{
	tree.EQL: true,
	tree.GT:  true,
	tree.GTE: true,
}

func (c *compiler) op(code uint16, k uint32) uint {
	ix := uint(len(c.result))
	c.result = append(c.result, unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
	return ix
}

func (c *compiler) moveAtoX() {
	c.op(A_TO_X, 0)
}

func (c *compiler) loadAt(pos uint32) {
	if c.currentlyLoaded != int(pos) {
		c.op(LOAD, pos)
		c.currentlyLoaded = int(pos)
	}
}

func (c *compiler) loadLiteral(lit uint32) {
	c.op(LOAD_VAL, lit)
	c.currentlyLoaded = -1
}

func (c *compiler) loadCurrentSyscall() {
	c.loadAt(syscallNameIndex)
}

func (c *compiler) performArithmetic(op tree.ArithmeticType, operand uint32) {
	switch op {
	case tree.PLUS:
		c.op(ADD_K, operand)
	case tree.MINUS:
		c.op(SUB_K, operand)
	case tree.MULT:
		c.op(MUL_K, operand)
	case tree.DIV:
		c.op(DIV_K, operand)
	case tree.BINAND:
		c.op(AND_K, operand)
	case tree.BINOR:
		c.op(OR_K, operand)
	case tree.LSH:
		c.op(LSH_K, operand)
	case tree.RSH:
		c.op(RSH_K, operand)
	}
}

func (c *compiler) checkCorrectSyscall(name string, setPosFlags bool) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	c.jumpOnKComparison(sys, tree.EQL, true, setPosFlags, false)
}

func (c *compiler) positiveAction(name string) {
	c.labelHere(positive)
	c.op(RET_K, SECCOMP_RET_ALLOW)
}

func (c *compiler) negativeAction(name string) {
	c.labelHere(negative)
	c.op(RET_K, SECCOMP_RET_KILL)
}
