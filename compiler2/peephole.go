package compiler2

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// This file contains space for implementing peephole optimization
// For now it doesn't do anything, but it could once we know what
// patterns show up

// One common pattern we likely will want to fix is:
// [ST n, LDX n]
// and rewrite it into [TAX]

// We might also see [LD_IMM v, ST n, LDX n]
// This should be rewritten into [LDX_IMM v]

// We should optimize [J* 0 1 *, JMP 1, JMP 1] to [J* 0 1]

// We will see a lot of things like [LD_IMM v, ST n, ... LDX n, <op>]
// These, both arithmetic and comparison ones should be rewritten to use the _K
// variants. However, these is easier done in the actual compiler at the moment
// These optimizations should also check if they operate on commutative
// operators and try to put the constant to the right, where K can be used.

func (c *compilerContext) optimizeCode() {
	index := 0

	// Do not pull out the length calculation here, since the length
	// can change during optimization
	for index < len(c.result) {
		c.optimizeAt(index)
		index++
	}
}

type optimizer func(*compilerContext, int) bool

var optimizers = []optimizer{
	doubleJumpOptimizer,
}

func (c *compilerContext) optimizeAt(i int) {
	for _, o := range optimizers {
		o(c, i)
	}
}

func isJump(s unix.SockFilter) bool {
	return s.Code&syscall.BPF_JMP != 0
}

func isUnconditionalJump(s unix.SockFilter) bool {
	return s.Code&(syscall.BPF_JMP|syscall.BPF_JA) != 0
}

func isConditionalJump(s unix.SockFilter) bool {
	return isJump(s) && !isUnconditionalJump(s)
}

func hasJumpTo(s unix.SockFilter, i int) bool {
	// TODO: continue
	return false
}

// doubleJumpOptimizer optimizes [J* 0 1 *, JMP 1, JMP 1] to [J* 0 1]
func doubleJumpOptimizer(c *compilerContext, i int) bool {
	// TODO: we should check all the conditions
	// TODO: we should make sure there are no other jumps coming in
	// TODO: we should make sure we fix up all other jumps when modifying the
	//       code

	// if len(c.result) > i+2 {
	// 	one, two, three := c.result[i], c.result[i+1], c.result[i+2]
	// 	if isConditionalJump(one) &&
	// 		isUnconditionalJump(two) &&
	// 		isUnconditionalJump(three) &&
	// 		hasJumpTo(one, i+1) &&
	// 		hasJumpTo(one, i+2) &&
	// 		hasUnconditionalJumpTo(two, i+2) &&
	// 		hasUnconditionalJumpTo(three, i+3) {
	// 		fmt.Printf("Hurrah\n")
	// 	}
	// }
	return false
}
