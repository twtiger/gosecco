package compiler2

import (
	"golang.org/x/sys/unix"
)

func (c *compilerContext) isNotLongJump(at, pos int) bool {
	return !((at-pos)-1 > c.maxJumpSize)
}

func (c *compilerContext) fixupJumps() {

	for l, at := range c.labels {
		for _, pos := range c.jts[l] {
			if c.isNotLongJump(at, pos) { // skip long jumps, we already fixed them up
				c.result[pos].Jt = uint8((at - pos) - 1)
			}
		}

		for _, pos := range c.jfs[l] {
			if c.isNotLongJump(at, pos) { // skip long jumps, we already fixed them up
				c.result[pos].Jf = uint8((at - pos) - 1)
			}
		}

		for _, pos := range c.uconds[l] {
			c.result[pos].K = uint32((at - pos) - 1)
		}
	}
}

func (c *compilerContext) hasPreviousUnconditionalJump(from int) bool {
	multJumps := false
	for _, v := range c.uconds {
		for _, pos := range v {
			if pos == from {
				multJumps = true
			}
		}
	}
	return multJumps
}

func (c *compilerContext) longJump(from int, positiveJump bool, to label) {
	hasPrev := c.hasPreviousUnconditionalJump(from)

	nextJ := from + 1
	if hasPrev {
		nextJ = from + 2
	}

	c.result = c.insertUnconditionalJump(nextJ)
	c.fixUpPreviousRule(from, positiveJump)
	c.shiftJumps(from, hasPrev)

	c.uconds[to] = append(c.uconds[to], nextJ)
}

func (c *compilerContext) insertUnconditionalJump(from int) []unix.SockFilter {
	var rules []unix.SockFilter
	x := unix.SockFilter{Code: OP_JMP_K, K: uint32(0)}

	for i, e := range c.result {
		if i == from {
			rules = append(rules, x)
		}
		rules = append(rules, e)
	}
	return rules
}

func shift(from int, incr int, elems map[label][]int) map[label][]int {
	jumps := make(map[label][]int, 0)

	for k, v := range elems {
		for _, pos := range v {
			if pos > from {
				pos += incr
			}
			jumps[k] = append(jumps[k], pos)
		}
	}
	return jumps
}

func shiftLabels(from int, incr int, elems map[label]int) map[label]int {
	labels := make(map[label]int, 0)

	for k, v := range elems {
		if v > from {
			v += incr
		}
		labels[k] = v
	}
	return labels
}

func (c *compilerContext) shiftJumps(from int, hasPrev bool) {
	incr := 1
	if hasPrev {
		incr = 2
	}

	c.jts = shift(from, incr, c.jts)
	c.jfs = shift(from, incr, c.jfs)
	c.uconds = shift(from, incr, c.uconds)
	c.labels = shiftLabels(from, incr, c.labels)
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
