package compiler

import (
	"github.com/twtiger/gosecco/tree"
)

type jumps string

const (
	jt      = "jt"
	jf      = "jf"
	neg     = "neg"
	chained = "chained"
)

func (c *compiler) positiveJumpTo(index uint, l label, neg, chained bool) {
	li := labelInfo{index, neg, chained}
	if l != noLabel {
		c.positiveLabels[l] = append(c.positiveLabels[l], li)
	}
}

func (c *compiler) negativeJumpTo(index uint, l label, neg, chained bool) {
	li := labelInfo{index, neg, chained}
	if l != noLabel {
		c.negativeLabels[l] = append(c.negativeLabels[l], li)
	}
}

func (c *compiler) jumpTo(num uint, terminalJF, terminalJT, neg, chained bool, jt, jf label) {
	if terminalJF {
		c.negativeJumpTo(num, jf, neg, chained)
	}
	if terminalJT {
		c.positiveJumpTo(num, jt, neg, chained)
	}
}

func (c *compiler) jumpOnKComparison(val uint32, cmp tree.ComparisonType, terminalJF, terminalJT, negated, chained bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, negated, chained, positive, negative)
	} else {
		c.jumpTo(num, terminalJT, terminalJF, negated, chained, negative, positive)
	}
}

func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, terminalJF, terminalJT, negated bool, chained bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, negated, chained, positive, negative)
	} else {
		c.jumpTo(num, terminalJT, terminalJF, negated, chained, negative, positive)
	}
}
