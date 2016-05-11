package compiler2

import (
	"golang.org/x/sys/unix"
)

func (c *compilerContext) fixupJumps() {

	for l, at := range c.labels {
		for _, pos := range c.jts[l] {
			if !((at-pos)-1 > c.maxJumpSize) { // skip long jumps, we already fixed them up
				c.result[pos].Jt = uint8((at - pos) - 1)
			}
		}

		for _, pos := range c.jfs[l] {
			if !((at-pos)-1 > c.maxJumpSize) { // skip long jumps, we already fixed them up
				c.result[pos].Jf = uint8((at - pos) - 1)
			}
		}

		// TODO: go through c.uconds and set K to be the correct value
	}
}

func (c *compilerContext) longJump(from int, positiveJump bool, to label) {

	//c.uconds[label] = append(c.uconds[label], from+1)
	// TODO to needs to be set in our unconditional jump list

	c.result = c.insertUnconditionalJump(from) // k needs to be set : from, to
	c.fixUpPreviousRule(from, positiveJump)
	c.shiftJumps(from)
}

func (c *compilerContext) insertUnconditionalJump(from int) []unix.SockFilter {
	rules := make([]unix.SockFilter, 0)
	at := len(c.result)
	k := uint32(at - from - 1)
	x := unix.SockFilter{Code: OP_JMP_K, K: k}

	rules = append(rules, c.result[:from+1]...)
	rules = append(rules, x)
	rules = append(rules, c.result[from+1:]...)
	return rules
}

func shift(from int, elems map[label][]int) map[label][]int {
	jumps := make(map[label][]int, 0)

	for k, v := range elems {
		for _, pos := range v {
			if pos >= from {
				pos += 1
			}
			jumps[k] = append(jumps[k], pos)
		}
	}
	return jumps
}

func shiftLabels(from int, elems map[label]int) map[label]int {
	labels := make(map[label]int, 0)
	for k, v := range elems {
		if v >= from {
			v += 1
		}
		labels[k] = v
	}
	return labels
}

func (c *compilerContext) shiftJumps(from int) {

	c.jts = shift(from, c.jts)
	c.jfs = shift(from, c.jfs)
	c.uconds = shift(from, c.uconds)
	c.labels = shiftLabels(from, c.labels)
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
