package compiler

type label string

type labelInfo struct {
	origin  uint
	negated bool
	chained bool
}

const (
	negative label = "negative"
	positive       = "positive"
	noLabel        = ""
)

func (c *compiler) labelHere(l label) {
	c.fixupJumpPoints(l, uint(len(c.result)))
}

func (c *compiler) fixupJumpPoints(l label, ix uint) {
	for _, e := range c.positiveLabels[l] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		if e.chained {
			c.result[e.origin].Jt = 0
		} else if e.negated {
			c.result[e.origin].Jf = uint8(ix-e.origin) - 1
		} else {
			c.result[e.origin].Jt = uint8(ix-e.origin) - 1
		}
	}
	delete(c.positiveLabels, l)

	for _, e := range c.negativeLabels[l] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		if e.chained {
			c.result[e.origin].Jf = 2
		} else if e.negated {
			c.result[e.origin].Jt = uint8(ix-e.origin) - 1
		} else {
			c.result[e.origin].Jf = uint8(ix-e.origin) - 1
		}
	}
	delete(c.negativeLabels, l)
}
