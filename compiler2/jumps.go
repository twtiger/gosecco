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
