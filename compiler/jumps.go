package compiler

import "golang.org/x/sys/unix"

type shift int

func (c *compilerContext) isLongJump(jumpSize int) bool {
	return jumpSize > c.maxJumpSize
}

func hasLongJump(index int, jts, jfs map[int]int) bool {
	// Using the unshifted index to look up positions in jts and jfs is
	// only safe if we're iterating backwards. Otherwise we would have to
	// fix up the positions in the maps as well and that would be fugly.

	if _, ok := jts[index]; ok {
		return true
	}
	if _, ok := jfs[index]; ok {
		return true
	}
	return false
}

func fixupWithShifts(pos, add int, shifts []shift) int {
	to := pos + add + 1
	currentAdd := add
	for _, s := range shifts {
		if int(s) > pos && int(s) <= to {
			currentAdd++
			to++
		}
	}
	return currentAdd
}

type longJumpContext struct {
	*compilerContext
	maxIndexWithLongJump     int
	jtLongJumps, jfLongJumps map[int]int
	shifts                   []shift
}

func (c *longJumpContext) fixupLongJumps() {
	// This is an optimization. Please don't comment away.
	if c.maxIndexWithLongJump == -1 {
		return
	}

	c.shifts = []shift{}

	currentIndex := c.maxIndexWithLongJump
	for currentIndex > -1 {
		current := c.result[currentIndex]

		if isConditionalJump(current) && hasLongJump(currentIndex, c.jtLongJumps, c.jfLongJumps) {
			hadJt := c.handleJTLongJumpFor(currentIndex)
			c.handleJFLongJumpFor(currentIndex, c.jfLongJumps, hadJt)
		} else {
			if isUnconditionalJump(current) {
				c.result[currentIndex].K = uint32(fixupWithShifts(currentIndex, int(c.result[currentIndex].K), c.shifts))
			} else {
				hadJt := c.shiftJt(currentIndex)
				c.shiftJf(hadJt, currentIndex)
			}
		}
		currentIndex--
	}
}

func (c *compilerContext) fixupJumps() {
	maxIndexWithLongJump := -1
	jtLongJumps := make(map[int]int)
	jfLongJumps := make(map[int]int)

	for l, at := range c.labels.allLabels() {
		for _, pos := range c.jts.allJumpsTo(l) {
			jumpSize := (at - pos) - 1
			if c.isLongJump(jumpSize) {
				if maxIndexWithLongJump < pos {
					maxIndexWithLongJump = pos
				}
				jtLongJumps[pos] = jumpSize
			} else {
				c.result[pos].Jt = uint8(jumpSize)
			}
		}

		for _, pos := range c.jfs.allJumpsTo(l) {
			jumpSize := (at - pos) - 1
			if c.isLongJump(jumpSize) {
				if maxIndexWithLongJump < pos {
					maxIndexWithLongJump = pos
				}
				jfLongJumps[pos] = jumpSize
			} else {
				c.result[pos].Jf = uint8(jumpSize)
			}
		}

		for _, pos := range c.uconds.allJumpsTo(l) {
			c.result[pos].K = uint32((at - pos) - 1)
		}
	}

	(&longJumpContext{c, maxIndexWithLongJump, jtLongJumps, jfLongJumps, nil}).fixupLongJumps()
}

func (c *longJumpContext) handleJTLongJumpFor(currentIndex int) bool {
	hadJt := false
	if jmpLen, ok := c.jtLongJumps[currentIndex]; ok {
		jmpLen = fixupWithShifts(currentIndex, jmpLen, c.shifts)
		hadJt = true

		newJf := int(c.result[currentIndex].Jf) + 1
		if c.isLongJump(newJf) {
			// Simple case, we can just add it to the long jumps for JF:
			c.jfLongJumps[currentIndex] = newJf
		} else {
			c.result[currentIndex].Jf = uint8(newJf)
		}

		c.insertJumps(currentIndex, jmpLen, 0)
	}
	return hadJt
}

func (c *longJumpContext) handleJFLongJumpFor(currentIndex int, jfLongJumps map[int]int, hadJt bool) {
	if jmpLen, ok := jfLongJumps[currentIndex]; ok {
		jmpLen = fixupWithShifts(currentIndex, jmpLen, c.shifts)
		var incr int
		incr, jmpLen = c.increment(hadJt, jmpLen, currentIndex)
		c.insertJumps(currentIndex, jmpLen, incr)
	}
}

func (c *longJumpContext) increment(hadJt bool, jmpLen, currentIndex int) (int, int) {
	incr := 0
	if hadJt {
		c.result[currentIndex+1].K++
		incr++
		jmpLen--
	} else {
		newJt := int(c.result[currentIndex].Jt) + 1
		if c.isLongJump(newJt) {
			// incr in this case doesn't seem to do much, all tests pass when it is changed to 0
			c.insertJumps(currentIndex, newJt, incr)
			incr++
		} else {
			c.result[currentIndex].Jt = uint8(newJt)
		}
	}
	return incr, jmpLen
}

func (c *longJumpContext) shiftJf(hadJt bool, currentIndex int) {
	newJf := fixupWithShifts(currentIndex, int(c.result[currentIndex].Jf), c.shifts)
	if c.isLongJump(newJf) {
		var incr int
		incr, _ = c.increment(hadJt, 0, currentIndex)
		c.insertJumps(currentIndex, newJf, incr)
	} else {
		c.result[currentIndex].Jf = uint8(newJf)
	}
}

func (c *longJumpContext) shiftJt(currentIndex int) bool {
	hadJt := false
	newJt := fixupWithShifts(currentIndex, int(c.result[currentIndex].Jt), c.shifts)
	if c.isLongJump(newJt) {
		hadJt = true

		// Jf doesn't need to be modified here, because it will be fixed up with the shifts. Hopefully correctly...
		c.insertJumps(currentIndex, newJt, 0)
	} else {
		c.result[currentIndex].Jt = uint8(newJt)
	}
	return hadJt
}

func (c *longJumpContext) insertJumps(currentIndex, pos, incr int) {
	c.insertUnconditionalJump(currentIndex+1+incr, pos)
	c.result[currentIndex].Jf = uint8(incr)
	c.shifts = append(c.shifts, shift(currentIndex+1+incr))
}

func insertSockFilter(sfs []unix.SockFilter, ix int, x unix.SockFilter) []unix.SockFilter {
	return append(
		append(
			append([]unix.SockFilter{}, sfs[:ix]...), x), sfs[ix:]...)
}

func (c *compilerContext) insertUnconditionalJump(from, k int) {
	x := unix.SockFilter{Code: OP_JMP_K, K: uint32(k)}
	c.result = insertSockFilter(c.result, from, x)
}

func (c *compilerContext) shiftJumpsBy(from, incr int) {
	c.jts.shift(from, incr)
	c.jfs.shift(from, incr)
	c.uconds.shift(from, incr)
	c.labels.shiftLabels(from, incr)
}
