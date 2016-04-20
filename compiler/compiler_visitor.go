package compiler

import (
	"github.com/twtiger/gosecco/tree"
)

type compilerVisitor struct {
	c *compiler
}

var compVals = map[tree.ComparisonType][]string{
	tree.EQL:  {"positive", "negative"},
	tree.GT:   {"positive", "negative"},
	tree.GTE:  {"positive", "negative"},
	tree.BIT:  {"positive", "negative"},
	tree.NEQL: {"negative", "positive"},
	tree.LT:   {"negative", "positive"},
	tree.LTE:  {"negative", "positive"},
}

func (cv *compilerVisitor) AcceptAnd(tree.And) {}

func (cv *compilerVisitor) AcceptArgument(a tree.Argument) {
	ix := ArgumentIndex[a.Index]
	cv.c.loadAt(ix["upper"])
	cv.c.jumpOnKComparison(0, tree.EQL, false, "positive", "negative")
	cv.c.loadAt(ix["lower"])
}

func (cv *compilerVisitor) AcceptArithmetic(a tree.Arithmetic) {
	a.Left.Accept(cv)
	rightOperand := a.Right.(tree.NumericLiteral)
	cv.c.performArithmetic(a.Op, rightOperand.Value)
}

func (cv *compilerVisitor) AcceptBinaryNegation(tree.BinaryNegation) {}
func (cv *compilerVisitor) AcceptBooleanLiteral(tree.BooleanLiteral) {}
func (cv *compilerVisitor) AcceptCall(tree.Call)                     {}

func (cv *compilerVisitor) AcceptComparison(c tree.Comparison) {
	act, _ := compVals[c.Op]
	lit, isLit := c.Right.(tree.NumericLiteral)
	if isLit {
		c.Left.Accept(cv)
		cv.c.jumpOnKComparison(lit.Value, c.Op, true, act[0], act[1])
	} else {
		c.Right.Accept(cv)
		cv.c.moveAtoX()
		c.Left.Accept(cv)
		cv.c.jumpOnXComparison(c.Op, act[0], act[1])
	}
}

func (cv *compilerVisitor) AcceptInclusion(tree.Inclusion) {}
func (cv *compilerVisitor) AcceptNegation(tree.Negation)   {}

func (cv *compilerVisitor) AcceptNumericLiteral(l tree.NumericLiteral) {
	cv.c.loadLiteral(l.Value)
}

func (cv *compilerVisitor) AcceptOr(tree.Or)             {}
func (cv *compilerVisitor) AcceptVariable(tree.Variable) {}

// func peepHole(filters []unix.SockFilter) []unix.SockFilter {
// 	one, two, three := filters[0], filters[1], filters[2]
// 	if one.Code == BPF_LD|BPF_IMM && two.Code == BPF_MISC|BPF_TAX && three.Code&(BPF_JMP|BPF_X) != 0 {
// 		return []unix.SockFilter{
// 			unix.SockFilter{
// 				Code: (three.Code & ^BPF_X) | BPF_K,
// 				Jt:   three.Jt,
// 				Jf:   three.Jf,
// 				K:    one.K,
// 			},
// 		}
// 	}
// 	return filters
// }
