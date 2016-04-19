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
	c.checkCorrectSyscall(r.Name)
	c.compileExpression(r.Body)
}

const syscallNameIndex = 0
const arg0IndexLowerWord = 0x10
const arg0IndexUpperWord = 0x14

var ComparisonOps = map[tree.ComparisonType]uint16{
	tree.EQL:  JEQ_K,
	tree.NEQL: JEQ_K,
	tree.GT:   JEG_K,
	tree.GTE:  JEGE_K,
	tree.LT:   JEG_K,
	tree.LTE:  JEGE_K,
	//BIT:  "bitSet",
}

const LOAD = BPF_LD | BPF_W | BPF_ABS
const LOAD_VAL = BPF_LD | BPF_IMM
const JEQ_K = BPF_JMP | BPF_JEQ | BPF_K
const JEG_K = BPF_JMP | BPF_JGT | BPF_K
const JEGE_K = BPF_JMP | BPF_JGE | BPF_K
const JEQ_X = BPF_JMP | BPF_JEQ | BPF_X
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

func (c *compiler) jumpOnComparison(val uint32, cmp tree.ComparisonType, jt, jf string) {
	jc := ComparisonOps[cmp]
	num := c.op(jc, val)
	c.positiveJumpTo(num, jt)
	c.negativeJumpTo(num, jf)
}

func (c *compiler) jumpIfEqualToX(jt, jf string) {
	num := c.op(JEQ_X, 0)
	c.positiveJumpTo(num, jt)
	c.negativeJumpTo(num, jf)
}

func (c *compiler) checkCorrectSyscall(name string) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	c.jumpOnComparison(sys, tree.EQL, "positive", "negative")
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
