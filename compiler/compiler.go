package compiler

import (
	"github.com/twtiger/gosecco/constants"
	"github.com/twtiger/gosecco/tree"
	"golang.org/x/sys/unix"
)

type label string

const (
	negative label = "negative"
	positive       = "positive"
	noLabel        = ""
)

func newCompiler() *compiler {
	return &compiler{
		currentlyLoaded: -1,
		positiveLabels:  make(map[label][]uint),
		negativeLabels:  make(map[label][]uint),
	}
}

func Compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c := newCompiler()
	c.compile(policy.Rules)
	return c.result, nil
}

type compiler struct {
	result          []unix.SockFilter
	currentlyLoaded int
	positiveLabels  map[label][]uint
	negativeLabels  map[label][]uint
}

func (c *compiler) compile(rules []tree.Rule) {
	for _, r := range rules {
		c.compileRule(r)
	}
	c.positiveAction(noLabel)
	c.negativeAction(noLabel)
}

func (c *compiler) compileExpression(x tree.Expression) {
	cv := &compilerVisitor{c, true, true}
	x.Accept(cv)
}

func (c *compiler) labelHere(l label) {
	c.fixupJumpPoints(l, uint(len(c.result)))
}

func (c *compiler) compileRule(r tree.Rule) {
	c.labelHere(negative)
	_, isBoolLit := r.Body.(tree.BooleanLiteral)
	c.checkCorrectSyscall(r.Name, isBoolLit) // set JT flag to final ret_allow only if the rule is a boolean literal
	c.compileExpression(r.Body)
}

const syscallNameIndex = 0

type argumentPosition struct {
	lower uint32
	upper uint32
}

var argument = []argumentPosition{
	argumentPosition{lower: 0x10, upper: 0x14},
	argumentPosition{lower: 0x18, upper: 0x1c},
	argumentPosition{lower: 0x20, upper: 0x24},
	argumentPosition{lower: 0x28, upper: 0x2c},
	argumentPosition{lower: 0x30, upper: 0x34},
	argumentPosition{lower: 0x38, upper: 0x3c},
}

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
	tree.BIT:  kexInstruction{k: JSET_K, x: JSET_X},
}

var posVals = map[tree.ComparisonType]bool{
	tree.EQL: true,
	tree.GT:  true,
	tree.GTE: true,
	tree.BIT: true,
}

const LOAD = BPF_LD | BPF_W | BPF_ABS
const LOAD_VAL = BPF_LD | BPF_IMM

const JEQ_K = BPF_JMP | BPF_JEQ | BPF_K
const JEQ_X = BPF_JMP | BPF_JEQ | BPF_X

const JEG_K = BPF_JMP | BPF_JGT | BPF_K
const JEG_X = BPF_JMP | BPF_JGT | BPF_X

const JEGE_K = BPF_JMP | BPF_JGE | BPF_K

const JEGE_X = BPF_JMP | BPF_JGE | BPF_X

const JSET_K = BPF_JMP | BPF_JSET | BPF_K
const JSET_X = BPF_JMP | BPF_JSET | BPF_X

const ADD_K = BPF_ALU | BPF_ADD | BPF_K
const MUL_K = BPF_ALU | BPF_MUL | BPF_K

const RET_K = BPF_RET | BPF_K
const A_TO_X = BPF_MISC | BPF_TAX

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

func (c *compiler) positiveJumpTo(index uint, l label) {
	if l != noLabel {
		c.positiveLabels[l] = append(c.positiveLabels[l], index)
	}
}

func (c *compiler) negativeJumpTo(index uint, l label) {
	if l != noLabel {
		c.negativeLabels[l] = append(c.negativeLabels[l], index)
	}
}

func (c *compiler) jumpTo(num uint, terminalJF, terminalJT bool, jt, jf label) {
	if terminalJF {
		c.negativeJumpTo(num, jf)
	}
	if terminalJT {
		c.positiveJumpTo(num, jt)
	}
}

func (c *compiler) jumpOnComparison(val uint32, cmp tree.ComparisonType) {
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	c.jumpTo(num, true, false, positive, negative)
}

func (c *compiler) jumpOnKComparison(val uint32, cmp tree.ComparisonType, terminalJF, terminalJT bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, positive, negative)
	} else {
		c.jumpTo(num, terminalJF, terminalJT, negative, positive)
	}
}

func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, terminalJF, terminalJT bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, positive, negative)
	} else {
		c.jumpTo(num, terminalJF, terminalJT, negative, positive)
	}
}

func (c *compiler) performArithmetic(op tree.ArithmeticType, operand uint32) {
	switch op {
	case tree.PLUS:
		c.op(ADD_K, operand)
	case tree.MULT:
		c.op(MUL_K, operand)
	}
}

func (c *compiler) checkCorrectSyscall(name string, setPosFlags bool) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	c.jumpOnKComparison(sys, tree.EQL, true, setPosFlags)
}

func (c *compiler) fixupJumpPoints(l label, ix uint) {
	for _, origin := range c.positiveLabels[l] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[origin].Jt = uint8(ix-origin) - 1
	}
	delete(c.positiveLabels, l)

	for _, origin := range c.negativeLabels[l] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[origin].Jf = uint8(ix-origin) - 1
	}
	delete(c.negativeLabels, l)
}

func (c *compiler) positiveAction(name string) {
	c.labelHere(positive)
	c.op(RET_K, SECCOMP_RET_ALLOW)
}

func (c *compiler) negativeAction(name string) {
	c.labelHere(negative)
	c.op(RET_K, SECCOMP_RET_KILL)
}
