package compiler

import "golang.org/x/sys/unix"

type shift struct {
	position, incr int
}

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
		if s.position > pos && s.position <= to {
			currentAdd++
			to++
		}
	}
	return currentAdd
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

	if maxIndexWithLongJump != -1 {
		shifts := []shift{}

		currentIndex := maxIndexWithLongJump
		for currentIndex > -1 {
			current := c.result[currentIndex]
			if isJump(current) {
				if isConditionalJump(current) &&
					hasLongJump(currentIndex, jtLongJumps, jfLongJumps) {

					hadPositive := false
					if jmpLen, ok := jtLongJumps[currentIndex]; ok {
						jmpLen = fixupWithShifts(currentIndex, jmpLen, shifts)

						hadPositive = true
						c.insertUnconditionalJump(currentIndex+1, jmpLen)
						c.result[currentIndex].Jt = 0

						newJf := int(c.result[currentIndex].Jf) + 1
						if c.isLongJump(newJf) {
							// Simple case, we can just add it to the long jumps for JF:
							jfLongJumps[currentIndex] = newJf
						} else {
							c.result[currentIndex].Jf = uint8(newJf)
						}

						shifts = append(shifts, shift{currentIndex + 1, 1})
					}

					if jmpLen, ok := jfLongJumps[currentIndex]; ok {
						jmpLen = fixupWithShifts(currentIndex, jmpLen, shifts)

						incr := 0
						if hadPositive {
							c.result[currentIndex+1].K++
							incr++
							jmpLen--
						} else {
							newJt := int(c.result[currentIndex].Jt) + 1
							if c.isLongJump(newJt) {
								c.insertUnconditionalJump(currentIndex+1, newJt)
								c.result[currentIndex].Jt = 0
								shifts = append(shifts, shift{currentIndex + 1, 1})
								incr++
							} else {
								c.result[currentIndex].Jt = uint8(newJt)
							}
						}
						c.insertUnconditionalJump(currentIndex+1+incr, jmpLen)
						c.result[currentIndex].Jf = uint8(incr)
						shifts = append(shifts, shift{currentIndex + 1 + incr, 1})
					}
				} else {
					if isUnconditionalJump(current) {
						c.result[currentIndex].K = uint32(fixupWithShifts(currentIndex, int(c.result[currentIndex].K), shifts))
					} else {
						hadPositive := false

						newJt := fixupWithShifts(currentIndex, int(c.result[currentIndex].Jt), shifts)
						if c.isLongJump(newJt) {
							hadPositive = true
							c.insertUnconditionalJump(currentIndex+1, newJt)
							c.result[currentIndex].Jt = 0

							// Jf doesn't need to be modified here, because it will be fixed up with the shifts. Hopefully correctly...

							shifts = append(shifts, shift{currentIndex + 1, 1})
						} else {
							c.result[currentIndex].Jt = uint8(newJt)
						}

						newJf := fixupWithShifts(currentIndex, int(c.result[currentIndex].Jf), shifts)
						if c.isLongJump(newJf) {
							incr := 0
							if hadPositive {
								c.result[currentIndex+1].K++
								incr++
							} else {
								newJt := int(c.result[currentIndex].Jt) + 1
								if c.isLongJump(newJt) {
									c.insertUnconditionalJump(currentIndex+1, newJt)
									c.result[currentIndex].Jt = 0
									shifts = append(shifts, shift{currentIndex + 1, 1})
									incr++
								} else {
									c.result[currentIndex].Jt = uint8(newJt)
								}
							}
							c.insertUnconditionalJump(currentIndex+1+incr, newJf)
							c.result[currentIndex].Jf = uint8(incr)
							shifts = append(shifts, shift{currentIndex + 1 + incr, 1})
						} else {
							c.result[currentIndex].Jf = uint8(newJf)
						}
					}
				}
			}

			currentIndex--
		}
	}
}

func (c *compilerContext) hasPreviousUnconditionalJump(from int) bool {
	return c.uconds.hasJumpFrom(from)
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
