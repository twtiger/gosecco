package compiler

import "github.com/twtiger/gosecco/tree"

func (c *compiler) positiveJumpTo(index uint, l label, neg bool) {
	li := labelInfo{index, neg}
	if l != noLabel {
		c.positiveLabels[l] = append(c.positiveLabels[l], li)
	}
}

func (c *compiler) negativeJumpTo(index uint, l label, neg bool) {
	li := labelInfo{index, neg}
	if l != noLabel {
		c.negativeLabels[l] = append(c.negativeLabels[l], li)
	}
}

func (c *compiler) jumpTo(num uint, terminalJF, terminalJT, neg bool, jt, jf label) {
	if terminalJF {
		c.negativeJumpTo(num, jf, neg)
	}
	if terminalJT {
		c.positiveJumpTo(num, jt, neg)
	}
}

func (c *compiler) jumpOnComparison(val uint32, cmp tree.ComparisonType) {
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	c.jumpTo(num, true, false, false, positive, negative)
}

func (c *compiler) jumpOnKComparison(val uint32, cmp tree.ComparisonType, terminalJF, terminalJT, negated bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, negated, positive, negative)
	} else {
		c.jumpTo(num, terminalJT, terminalJF, negated, negative, positive)
	}
}

func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, terminalJF, terminalJT, negated bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, negated, positive, negative)
	} else {
		c.jumpTo(num, terminalJT, terminalJF, negated, negative, positive)
	}
}
