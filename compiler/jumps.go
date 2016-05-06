package compiler

import (
	"github.com/twtiger/gosecco/tree"
)

type jumps string

const (
	jt = "jt"
	jf = "jf"
)

func (c *compiler) positiveJumpTo(index uint, l label) {
	if l != noLabel {
		c.positiveLabels[l] = append(c.positiveLabels[l], index)
	}
}

func (c *compiler) negativeJumpTo(index uint, l label) {
	if l != noLabel {
		c.negativeLabels[l] = append(c.negativeLabels[l], index)
	}
}

func (c *compiler) jumpTo(num uint, jt, jf label) {
	c.positiveJumpTo(num, jt)
	c.negativeJumpTo(num, jf)
}

func (c *compiler) jumpOnKComp(val uint32, cmp tree.ComparisonType, jt, jf label) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if isPos {
		c.jumpTo(num, jt, jf)
	} else {
		c.jumpTo(num, jf, jt)
	}
}

func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, jt, jf label) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)

	if isPos {
		c.jumpTo(num, jt, jf)
	} else {
		c.jumpTo(num, jf, jt)
	}
}
