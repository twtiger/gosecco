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

type jumpType string

const (
	TermJf  = "TermJf"
	TermJ   = "TermJ"
	ChainJ  = "ChainJ"
	ChainJt = "ChainJt"
	ExlHi   = "ExlHi"
)

type jumpPoint struct {
	jf, jt, chained bool
}

// need to take another look at the chain attribute
var jumpPoints = map[jumpType]jumpPoint{
	TermJf:  jumpPoint{jf: true, jt: false, chained: false},
	TermJ:   jumpPoint{jf: true, jt: true, chained: false},
	ChainJ:  jumpPoint{jf: true, jt: true, chained: true},
	ExlHi:   jumpPoint{jf: true, jt: false, chained: false},
	ChainJt: jumpPoint{jf: false, jt: true, chained: false},
}

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

func (c *compiler) jumpOnSyscallComparison(val uint32, cmp tree.ComparisonType, terminalJF, terminalJT bool) {
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	c.jumpTo(num, terminalJF, terminalJT, false, false, positive, negative)
}

func (c *compiler) jumpOnKComp(val uint32, cmp tree.ComparisonType, jp jumpPoint, negated, inverted bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].k
	num := c.op(jc, val)
	if !isPos || negated {
		c.jumpTo(num, jp.jt, jp.jf, inverted, jp.chained, negative, positive)
	} else {
		c.jumpTo(num, jp.jf, jp.jt, inverted, jp.chained, positive, negative)
	}
}

//TODO add negation here
func (c *compiler) jumpOnXComparison(cmp tree.ComparisonType, jp jumpPoint, inverted bool) {
	_, isPos := posVals[cmp]
	jc := comparisonOps[cmp].x
	num := c.op(jc, 0)
	if isPos {
		c.jumpTo(num, jp.jf, jp.jt, inverted, jp.chained, positive, negative)
	} else {
		c.jumpTo(num, jp.jt, jp.jf, inverted, jp.chained, negative, positive)
	}
}
