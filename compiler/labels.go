package compiler

type label string

const (
	negative label = "negative"
	positive       = "positive"
	next           = "next"
	noLabel        = ""
)

func (c *compiler) labelHere(l label) {
	c.fixupJumpPoints(l, uint(len(c.result)))
}

func (c *compiler) fixupJumpPoints(l label, ix uint) {
	for _, e := range c.positiveLabels[l] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[e].Jt = uint8(ix-e) - 1
	}
	delete(c.positiveLabels, l)

	for _, e := range c.negativeLabels[l] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[e].Jf = uint8(ix-e) - 1
	}
	delete(c.negativeLabels, l)
}
