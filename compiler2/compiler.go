package compiler2

import (
	"fmt"

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
// TODO: check that the stack doesn't overflow

type label string

type compilerContext struct {
	result       []unix.SockFilter
	stackTop     uint32
	jts          map[label][]int
	jfs          map[label][]int
	labels       map[label]int
	labelCounter int
}

func createCompilerContext() *compilerContext {
	return &compilerContext{
		jts:    make(map[label][]int),
		jfs:    make(map[label][]int),
		labels: make(map[label]int),
	}
}

func (c *compilerContext) op(code uint16, k uint32) {
	c.result = append(c.result, unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
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

// TODO: check we are not outside limits here
func (c *compilerContext) pushAToStack() error {
	c.op(OP_STORE, c.stackTop)
	c.stackTop++
	return nil
}

func (c *compilerContext) popStackToX() error {
	c.stackTop--
	c.op(OP_LOAD_MEM_X, c.stackTop)
	return nil
}
