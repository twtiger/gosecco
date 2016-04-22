package compiler

import (
	"github.com/twtiger/gosecco/constants"
	"github.com/twtiger/gosecco/tree"
	"golang.org/x/sys/unix"
)

func newCompiler() *compiler {
	return &compiler{
		currentlyLoaded: -1,
		positiveLabels:  make(map[string][]uint),
		negativeLabels:  make(map[string][]uint),
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
	positiveLabels  map[string][]uint
	negativeLabels  map[string][]uint
}

func (c *compiler) compile(rules []tree.Rule) {
	for _, r := range rules {
		c.compileRule(r)
	}
	c.positiveAction("")
	c.negativeAction("")
}

func (c *compiler) compileExpression(x tree.Expression) {
	cv := &compilerVisitor{c, true}
	x.Accept(cv)
}

func (c *compiler) labelHere(label string) {
	c.fixupJumpPoints(label, uint(len(c.result)))
}

func (c *compiler) compileRule(r tree.Rule) {
	c.labelHere("negative")
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

func (c *compiler) positiveJumpTo(index uint, label string) {
	if label != "" {
		c.positiveLabels[label] = append(c.positiveLabels[label], index)
	}
}

func (c *compiler) negativeJumpTo(index uint, label string) {
	if label != "" {
		c.negativeLabels[label] = append(c.negativeLabels[label], index)
	}
}

func (c *compiler) jumpOnKComparison(val uint32, cmp tree.ComparisonType, setPosFlags, isTerminal bool, jt, jf string) {
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if setPosFlags {
		c.positiveJumpTo(num, jt)
	}
	if isTerminal {
		c.negativeJumpTo(num, jf)
	}
}

func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, isTerminal bool, jt, jf string) {
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)
	c.positiveJumpTo(num, jt)
	if isTerminal {
		c.negativeJumpTo(num, jf)
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
	c.jumpOnKComparison(sys, tree.EQL, setPosFlags, true, "positive", "negative")
}

func (c *compiler) fixupJumpPoints(label string, ix uint) {
	for _, origin := range c.positiveLabels[label] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[origin].Jt = uint8(ix-origin) - 1
	}
	delete(c.positiveLabels, label)

	for _, origin := range c.negativeLabels[label] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[origin].Jf = uint8(ix-origin) - 1
	}
	delete(c.negativeLabels, label)
}

func (c *compiler) positiveAction(name string) {
	c.labelHere("positive")
	c.op(RET_K, SECCOMP_RET_ALLOW)
}

func (c *compiler) negativeAction(name string) {
	c.labelHere("negative")
	c.op(RET_K, SECCOMP_RET_KILL)
}
