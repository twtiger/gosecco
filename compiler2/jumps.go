package compiler2

func (c *compilerContext) fixupJumps() {
	// We can ignore long jumps here - they should have been fixed when the label was registered - easier that way
	for l, at := range c.labels {
		for _, pos := range c.jts[l] {
			c.result[pos].Jt = uint8((at - pos) - 1)
		}
		for _, pos := range c.jfs[l] {
			c.result[pos].Jf = uint8((at - pos) - 1)
		}
	}
}

// longJump will take the information necessary to make a long jump
// from is the index of the instruction that the long jump was supposed to be from
// positiveJump is whether we are talking jt or jf
// to is the label to jump to
// This function will shift everything necessary in things that keep track of other jumps
func (c *compilerContext) longJump(from uint16, positiveJump bool, to label) {
}
