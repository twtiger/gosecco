package compiler

import (
	"golang.org/x/sys/unix"
)

func (c *compilerContext) isNotLongJump(at, pos int) bool {
	return !((at-pos)-1 > c.maxJumpSize)
}

func (c *compilerContext) fixupJumps() {
	for l, at := range c.labels.allLabels() {
		for _, pos := range c.jts.allJumpsTo(l) {
			if c.isNotLongJump(at, pos) { // skip long jumps, we already fixed them up
				c.result[pos].Jt = uint8((at - pos) - 1)
			}
		}

		for _, pos := range c.jfs.allJumpsTo(l) {
			if c.isNotLongJump(at, pos) { // skip long jumps, we already fixed them up
				c.result[pos].Jf = uint8((at - pos) - 1)
			}
		}

		for _, pos := range c.uconds.allJumpsTo(l) {
			c.result[pos].K = uint32((at - pos) - 1)
		}
	}
}

func (c *compilerContext) hasPreviousUnconditionalJump(from int) bool {
	return c.uconds.hasJumpFrom(from)
}

func (c *compilerContext) longJump(from int, positiveJump bool, to label) {
	hasPrev := c.hasPreviousUnconditionalJump(from)

	nextJ := from + 1
	if hasPrev {
		nextJ = from + 2
	}

	c.insertUnconditionalJump(nextJ)
	c.fixUpPreviousRule(from, positiveJump)
	c.shiftJumps(from, hasPrev)

	c.uconds.registerJump(to, nextJ)
}

func insertSockFilter(sfs []unix.SockFilter, ix int, x unix.SockFilter) []unix.SockFilter {
	return append(
		append(
			append([]unix.SockFilter{}, sfs[:ix]...), x), sfs[ix:]...)
}

func (c *compilerContext) insertUnconditionalJump(from int) {
	x := unix.SockFilter{Code: OP_JMP_K, K: uint32(0)}
	c.result = insertSockFilter(c.result, from, x)
}

func (c *compilerContext) shiftJumps(from int, hasPrev bool) {
	incr := 1
	if hasPrev {
		incr = 2
	}
	c.shiftJumpsBy(from, incr)
}

func (c *compilerContext) shiftJumpsBy(from, incr int) {
	c.jts.shift(from, incr)
	c.jfs.shift(from, incr)
	c.uconds.shift(from, incr)
	c.labels.shiftLabels(from, incr)
}

func (c *compilerContext) fixUpPreviousRule(from int, positiveJump bool) {
	if positiveJump {
		c.result[from].Jt = 0
		c.result[from].Jf = 1
	} else {
		c.result[from].Jt = 1
		c.result[from].Jf = 0
	}
}
