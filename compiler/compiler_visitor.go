package compiler

import (
	"github.com/twtiger/gosecco/tree"
)

type compilerVisitor struct {
	c          *compiler
	terminalJF bool
	terminalJT bool
}

func (cv *compilerVisitor) AcceptArgument(a tree.Argument) {
	ix := argument[a.Index]
	cv.c.loadAt(ix.upper)
	cv.c.jumpOnComparison(0, tree.EQL)
	cv.c.loadAt(ix.lower)
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
	lit, isLit := c.Right.(tree.NumericLiteral)
	if isLit {
		c.Left.Accept(cv)
		cv.c.jumpOnKComparison(lit.Value, c.Op, cv.terminalJF, cv.terminalJT)
	} else {
		c.Right.Accept(cv)
		cv.c.moveAtoX()
		c.Left.Accept(cv)
		cv.c.jumpOnXComparison(c.Op, cv.terminalJF, cv.terminalJT)
	}
}

func (cv *compilerVisitor) AcceptInclusion(tree.Inclusion) {}
func (cv *compilerVisitor) AcceptNegation(tree.Negation)   {}

func (cv *compilerVisitor) AcceptNumericLiteral(l tree.NumericLiteral) {
	cv.c.loadLiteral(l.Value)
}

func (cv *compilerVisitor) AcceptAnd(c tree.And) {
	cv.terminalJT = false
	c.Left.Accept(cv)
	cv.terminalJT = true
	c.Right.Accept(cv)
}

func (cv *compilerVisitor) AcceptOr(c tree.Or) {
	cv.terminalJF = false
	c.Left.Accept(cv)
	cv.terminalJF = true
	c.Right.Accept(cv)
}
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
