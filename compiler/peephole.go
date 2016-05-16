package compiler

import (
	"fmt"
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

// Some patterns look amenable to optimization but in practice won't be
// - it's important that we are wary of trying to fix up jumps too much.

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
		if !c.optimizeAt(index) {
			index++
		}
	}
}

type optimizer func(*compilerContext, int) bool

var optimizers = []optimizer{
	jumpAfterConditionalJumpOptimizer,
}

func (c *compilerContext) optimizeAt(i int) bool {
	optimized := false
	for _, o := range optimizers {
		if o(c, i) {
			optimized = true
		}
	}
	return optimized
}

func isJump(s unix.SockFilter) bool {
	return bpfClass(s.Code) == syscall.BPF_JMP
}

func isConditionalJump(s unix.SockFilter) bool {
	return isJump(s) && bpfOp(s.Code) != syscall.BPF_JA
}

func isUnconditionalJump(s unix.SockFilter) bool {
	return isJump(s) && bpfOp(s.Code) == syscall.BPF_JA
}

// hasJumpTarget will return true if the conditional jump given has at least one of
// its target being the potential target.
func hasJumpTarget(c *compilerContext, conditionalIndex, potentialTarget int) bool {
	return c.jts.hasJumpTarget(c, conditionalIndex, potentialTarget) ||
		c.jfs.hasJumpTarget(c, conditionalIndex, potentialTarget)
}

// hasOnlyJumpFrom will make sure that the given jump target is only jumped to from the
// given expected from location - this will return false if the expectedFrom is a
// conditional jump where both conditions point to the jump target
func hasOnlyJumpFrom(c *compilerContext, jumpTarget, expectedFrom int) bool {
	return (c.jts.countJumpsFrom(c, jumpTarget, expectedFrom)+
		c.jfs.countJumpsFrom(c, jumpTarget, expectedFrom)+
		c.uconds.countJumpsFrom(c, jumpTarget, expectedFrom)) == 1 &&
		(c.jts.countJumpsFromAny(c, jumpTarget)+
			c.jfs.countJumpsFromAny(c, jumpTarget)+
			c.uconds.countJumpsFromAny(c, jumpTarget)) == 1
}

// isNotOversizedJump takes a pointer to an unconditional jump and return true if
// the jump is smaller than the max jump size.
func isNotOversizedJump(c *compilerContext, jumpPoint int) bool {
	return !c.uconds.jumpSizeIsOversized(c, jumpPoint)
}

func redirectJumpOf(c *compilerContext, ucond, cond int) {
	newJumpTarget := c.uconds.jumpTargetOf(ucond)
	c.uconds.removeJumpTarget(ucond)
	var sourceJm *jumpMap

	if c.jts.countJumpsFrom(c, ucond, cond) == 1 {
		sourceJm = c.jts
	} else if c.jfs.countJumpsFrom(c, ucond, cond) == 1 {
		sourceJm = c.jfs
	} else {
		panic(fmt.Sprintf("No jumps to redirect (programmer error): ucond: %d cond: %d\n%#v\n%#v\n", ucond, cond, c.jts, c.labels))
	}

	oldLabel := c.labels.labelsAt(ucond)[0]
	sourceJm.redirectJump(oldLabel, newJumpTarget)
	c.labels.removeLabel(oldLabel)
}

func (c *compilerContext) removeInstructionAt(index int) {
	c.result = append(c.result[:index], c.result[index+1:]...)
}

// jumpAfterConditionalJumpOptimizer will optimize situations where a JMP instruction
// directly follows a conditional jump where one of the arms of the conditional jump
// is zero. It will make sure that no other jump points end up on the specific JMP instruction
// before removing it. It will also make sure the resulting jump is not too large.
// An example of a fragment that would be changed would be this:
//    jeq_k	00	01	3D
//    jmp	13
// This can be optimized to:
//    jeq_k	13	00	3D
func jumpAfterConditionalJumpOptimizer(c *compilerContext, ix int) bool {
	optimized := false

	if ix+1 < len(c.result) {
		oneIndex, twoIndex := ix, ix+1
		one, two := c.result[oneIndex], c.result[twoIndex]
		if isConditionalJump(one) &&
			isUnconditionalJump(two) &&
			hasJumpTarget(c, oneIndex, twoIndex) &&
			hasOnlyJumpFrom(c, twoIndex, oneIndex) &&
			isNotOversizedJump(c, twoIndex) {

			redirectJumpOf(c, twoIndex, oneIndex)
			c.shiftJumpsBy(oneIndex+1, -1)
			c.removeInstructionAt(twoIndex)

			optimized = true
		}
	}

	return optimized
}
