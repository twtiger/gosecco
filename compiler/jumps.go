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

func (c *compiler) negativeJumpTo(index uint, l label, inv, chained bool) {
	li := labelInfo{index, inv, chained}
	if l != noLabel {
		c.negativeLabels[l] = append(c.negativeLabels[l], li)
	}
}

func (c *compiler) jumpTo(num uint, terminalJF, terminalJT, inv, chained bool, jt, jf label) {
	if terminalJF {
		c.negativeJumpTo(num, jf, inv, chained)
	}
	if terminalJT {
		c.positiveJumpTo(num, jt, inv, chained)
	}
}

func (c *compiler) jumpOnKComparison(val uint32, cmp tree.ComparisonType, terminalJF, terminalJT, negated, inverted, chained bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if !isPos || negated {
		c.jumpTo(num, terminalJT, terminalJF, inverted, chained, negative, positive)
	} else {
		c.jumpTo(num, terminalJF, terminalJT, inverted, chained, positive, negative)
	}
}

//TODO add negation here
func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, terminalJF, terminalJT, inverted bool, chained bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)
	if isPos {
		c.jumpTo(num, terminalJF, terminalJT, inverted, chained, positive, negative)
	} else {
		c.jumpTo(num, terminalJT, terminalJF, inverted, chained, negative, positive)
	}
}
