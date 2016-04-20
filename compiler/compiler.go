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
	cv := &compilerVisitor{c}
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
const arg0IndexLowerWord = 0x10
const arg0IndexUpperWord = 0x14

var ComparisonOps = map[tree.ComparisonType]map[string]uint16{
	tree.EQL:  {"K": JEQ_K, "X": JEQ_X},
	tree.NEQL: {"K": JEQ_K, "X": JEQ_X},
	tree.GT:   {"K": JEG_K, "X": JEG_X},
	tree.GTE:  {"K": JEGE_K, "X": JEGE_X},
	tree.LT:   {"K": JEG_K, "X": JEG_X},
	tree.LTE:  {"K": JEGE_K, "X": JEGE_X},
	tree.BIT:  {"K": JSET_K, "X": JSET_X},
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

func (c *compiler) jumpOnKComparison(val uint32, cmp tree.ComparisonType, setPosFlags bool, jt, jf string) {
	jc := ComparisonOps[cmp]["K"]
	num := c.op(jc, val)
	if setPosFlags {
		c.positiveJumpTo(num, jt)
	}
	c.negativeJumpTo(num, jf)
}

func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, jt, jf string) {
	jc := ComparisonOps[cmp]["X"]
	num := c.op(jc, 0)
	c.positiveJumpTo(num, jt)
	c.negativeJumpTo(num, jf)
}

func (c *compiler) performArithmetic(op tree.ArithmeticType, operand uint32) {
	c.op(ADD_K, operand)
}

func (c *compiler) checkCorrectSyscall(name string, setPosFlags bool) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	c.jumpOnKComparison(sys, tree.EQL, setPosFlags, "positive", "negative")
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
