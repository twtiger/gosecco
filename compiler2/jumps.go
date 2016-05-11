package compiler2

import (
	//	"fmt"
	"golang.org/x/sys/unix"
)

func (c *compilerContext) fixupJumps() {

	//	fmt.Println("jump trues", c.jts)
	//	fmt.Println("jump false", c.jfs)
	//	fmt.Println("unconditional jumps", c.uconds)
	//	fmt.Println("labels", c.labels)
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

		for _, pos := range c.uconds[l] {
			c.result[pos].K = uint32((at - pos) - 1)
		}
	}
}

func (c *compilerContext) longJump(from int, positiveJump bool, to label) {
	c.result = c.insertUnconditionalJump(from) // k needs to be set : from, to
	c.fixUpPreviousRule(from, positiveJump)
	c.shiftJumps(from)
	c.uconds[to] = append(c.uconds[to], from+1)
}

func (c *compilerContext) insertUnconditionalJump(from int) []unix.SockFilter {
	rules := make([]unix.SockFilter, 0)
	k := uint32(0)
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
