package compiler

import (
	"fmt"

	"github.com/twtiger/gosecco/tree"
)

type compilerVisitor struct {
	c          *compiler
	terminalJF bool
	terminalJT bool
	negated    bool
	topLevel   bool
}

func (cv *compilerVisitor) AcceptArgument(a tree.Argument) {
	cv.topLevel = false
	ix := argument[a.Index]
	cv.c.loadAt(ix.upper)
	cv.c.jumpOnComparison(0, tree.EQL)
	cv.c.loadAt(ix.lower)
}

func (cv *compilerVisitor) AcceptArithmetic(a tree.Arithmetic) {
	cv.topLevel = false
	a.Left.Accept(cv)
	rightOperand := a.Right.(tree.NumericLiteral)
	cv.c.performArithmetic(a.Op, rightOperand.Value)
}

func (cv *compilerVisitor) AcceptBinaryNegation(tree.BinaryNegation) {
	cv.topLevel = false
}

func (cv *compilerVisitor) AcceptBooleanLiteral(val tree.BooleanLiteral) {
	if cv.topLevel {
		// TODO: compile here
	} else {
		panic(fmt.Sprintf("Programming error: there should never be any boolean literals left outside of the toplevel if the simplifier works correctly: syscall: %s - %s", cv.c.currentlyCompilingSyscall, tree.ExpressionString(cv.c.currentlyCompilingExpression)))
	}
	cv.topLevel = false
}

func (cv *compilerVisitor) AcceptCall(tree.Call) {
	panic(fmt.Sprintf("Programming error: there should never be any unexpanded calls if the unifier works correctly: syscall: %s - %s", cv.c.currentlyCompilingSyscall, tree.ExpressionString(cv.c.currentlyCompilingExpression)))
}

func (cv *compilerVisitor) AcceptComparison(c tree.Comparison) {
	cv.topLevel = false
	lit, isLit := c.Right.(tree.NumericLiteral)
	if isLit {
		c.Left.Accept(cv)
		cv.c.jumpOnKComparison(lit.Value, c.Op, cv.terminalJF, cv.terminalJT, cv.negated)
	} else {
		c.Right.Accept(cv)
		cv.c.moveAtoX()
		c.Left.Accept(cv)
		cv.c.jumpOnXComparison(c.Op, cv.terminalJF, cv.terminalJT, cv.negated)
	}
}

func (cv *compilerVisitor) toggleTerminalJumps(b bool) {
	cv.topLevel = false
	if b == true {
		cv.terminalJF = !cv.terminalJF
	} else {
		cv.terminalJT = !cv.terminalJT
	}
}

func (cv *compilerVisitor) AcceptInclusion(c tree.Inclusion) {
	cv.topLevel = false
	if c.Positive == false {
		cv.negated = true
	}
	c.Left.Accept(cv)
	cv.toggleTerminalJumps(c.Positive)

	_, isLit := c.Left.(tree.NumericLiteral)
	if isLit {
		cv.c.moveAtoX()

		for i, e := range c.Rights {
			if i == len(c.Rights)-1 {
				cv.toggleTerminalJumps(c.Positive)
			}
			e.Accept(cv)
			cv.c.jumpOnXComparison(tree.EQL, cv.terminalJF, cv.terminalJT, cv.negated)
		}
	} else {

		for i, e := range c.Rights {
			if i == len(c.Rights)-1 {
				cv.toggleTerminalJumps(c.Positive)
			}
			lit, isLiteral := e.(tree.NumericLiteral)
			if isLiteral {
				cv.c.jumpOnKComparison(lit.Value, tree.EQL, cv.terminalJF, cv.terminalJT, cv.negated)
			} else {
				cv.c.moveAtoX()
				e.Accept(cv)
				cv.c.jumpOnXComparison(tree.EQL, cv.terminalJF, cv.terminalJT, cv.negated)
			}
		}
	}
}

func (cv *compilerVisitor) AcceptNegation(c tree.Negation) {
	cv.topLevel = false
	cv.negated = true
	c.Operand.Accept(cv)
}

func (cv *compilerVisitor) AcceptNumericLiteral(l tree.NumericLiteral) {
	cv.topLevel = false
	cv.c.loadLiteral(l.Value)
}

func (cv *compilerVisitor) AcceptAnd(c tree.And) {
	cv.topLevel = false
	cv.terminalJT = !cv.terminalJT
	c.Left.Accept(cv)
	cv.terminalJT = !cv.terminalJT
	c.Right.Accept(cv)
}

func (cv *compilerVisitor) AcceptOr(c tree.Or) {
	cv.topLevel = false
	cv.terminalJF = false
	c.Left.Accept(cv)
	cv.terminalJF = true
	c.Right.Accept(cv)
}

func (cv *compilerVisitor) AcceptVariable(tree.Variable) {
	panic(fmt.Sprintf("Programming error: there should never be any unexpanded variables if the unifier works correctly: syscall: %s - %s", cv.c.currentlyCompilingSyscall, tree.ExpressionString(cv.c.currentlyCompilingExpression)))
}

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
