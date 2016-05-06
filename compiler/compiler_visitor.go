package compiler

import (
	"fmt"

	"github.com/twtiger/gosecco/tree"
)

type compilerVisitor struct {
	c        *compiler
	topLevel bool
	jf, jt   label
}

func getLower(k uint64) uint32 {
	return uint32(k)
}

func getUpper(k uint64) uint32 {
	return uint32(k >> 32)
}

func (cv *compilerVisitor) AcceptArgument(a tree.Argument) {
	cv.topLevel = false
	ix := argument[a.Index]
	switch a.Type {
	case tree.Hi:
		cv.c.loadAt(ix.upper)
	case tree.Low:
		cv.c.loadAt(ix.lower)
	default:
		panic(fmt.Sprintf("Incorrect argument type"))
	}
}

func (cv *compilerVisitor) AcceptArithmetic(a tree.Arithmetic) {
	cv.topLevel = false
	a.Left.Accept(cv)
	rightOperand := a.Right.(tree.NumericLiteral)
	cv.c.performArithmetic(a.Op, uint32(rightOperand.Value))
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

func detectSpecialCasesOn(e tree.Expression) (*tree.Argument, *tree.NumericLiteral, bool, bool) {
	switch et := e.(type) {
	case tree.Argument:
		if et.Type == tree.Full {
			return &et, nil, true, false
		}
	case tree.NumericLiteral:
		return nil, &et, false, true
	}
	return nil, nil, false, false
}

func detectSpecialCases(c tree.Comparison) (argL *tree.Argument, argR *tree.Argument, litL *tree.NumericLiteral, litR *tree.NumericLiteral, leftIsArg bool, rightIsArg bool, leftIsLit bool, rightIsLit bool) {
	argL, litL, leftIsArg, leftIsLit = detectSpecialCasesOn(c.Left)
	argR, litR, rightIsArg, rightIsLit = detectSpecialCasesOn(c.Right)
	return
}

func (cv *compilerVisitor) compareExpressionToArg(a *tree.Argument, e tree.Expression, op tree.ComparisonType) {
	e.Accept(cv)
	cv.c.moveAtoX()
	lx := argument[a.Index]
	cv.c.loadAt(lx.upper)

	switch op {
	case tree.NEQL:
		cv.c.jumpOnXComparison(op, cv.jt, noLabel)
	case tree.EQL:
		cv.c.jumpOnXComparison(op, next, cv.jf)
	}

	cv.c.loadAt(lx.lower)
	cv.c.jumpOnXComparison(op, cv.jt, cv.jf)
}

func (cv *compilerVisitor) AcceptComparison(c tree.Comparison) {
	cv.topLevel = false
	argL, argR, litL, litR, leftArg, rightArg, leftLit, rightLit := detectSpecialCases(c)

	if leftArg && rightLit {
		ix := argument[argL.Index]
		cv.jumpOnK(litR.Value, ix, c.Op)
	}

	if leftLit && rightArg {
		ix := argument[argR.Index]
		cv.jumpOnK(litL.Value, ix, c.Op)
	}

	if leftArg && rightArg {
		rx := argument[argR.Index]
		lx := argument[argL.Index]
		cv.jumpOnX(rx, lx, c.Op)
	}

	if !rightArg && !rightLit && leftArg {
		cv.compareExpressionToArg(argL, c.Right, c.Op)
	}

	if !leftArg && !leftLit && rightArg {
		cv.compareExpressionToArg(argR, c.Left, c.Op)
	}

	if !leftLit && !leftArg && !rightLit && !rightArg {
		c.Left.Accept(cv)
		cv.c.moveAtoX()
		c.Right.Accept(cv)
		cv.c.jumpOnXComparison(c.Op, cv.jt, cv.jf)
	}

}

var count = 0

func nextLabel() label {
	count += 1
	return label(fmt.Sprintf("%d", count))
}

func (cv *compilerVisitor) jumpOnK(l uint64, ix argumentPosition, op tree.ComparisonType) {
	cv.c.loadAt(ix.upper)
	next := nextLabel()

	switch op {
	case tree.NEQL:
		cv.c.jumpOnKComp(getUpper(l), op, cv.jt, noLabel)
	case tree.EQL:
		cv.c.jumpOnKComp(getUpper(l), op, next, cv.jf)
	case tree.GT:
		cv.c.jumpOnKComp(getUpper(l), op, next, cv.jf)
	case tree.LT:
		cv.c.jumpOnKComp(getUpper(l), op, next, cv.jf)
	}
	cv.c.labelHere(next)

	cv.c.loadAt(ix.lower)
	cv.c.jumpOnKComp(getLower(l), op, cv.jt, cv.jf)
}

func (cv *compilerVisitor) jumpOnX(ix argumentPosition, rx argumentPosition, op tree.ComparisonType) {
	cv.c.loadAt(ix.upper)
	cv.c.moveAtoX()
	cv.c.loadAt(rx.upper)
	next := nextLabel()

	switch op {
	case tree.NEQL:
		cv.c.jumpOnXComparison(op, cv.jt, noLabel)
	case tree.EQL:
		cv.c.jumpOnXComparison(op, next, cv.jf)
	}
	cv.c.labelHere(next)

	cv.c.loadAt(ix.lower)
	cv.c.moveAtoX()
	cv.c.loadAt(rx.lower)
	cv.c.jumpOnXComparison(op, cv.jt, cv.jf)
}

func (cv *compilerVisitor) setJumpPoints(p bool) {
	if !p {
		cv.jt = negative
		cv.jf = positive
	} else {
		cv.jt = positive
		cv.jf = negative
	}
}

func (cv *compilerVisitor) goToNextComparison(isExclusive bool) label {
	n := nextLabel()
	if isExclusive {
		cv.jt = n
	} else {
		cv.jf = n
	}
	return n
}

func (cv *compilerVisitor) AcceptInclusion(c tree.Inclusion) {
	cv.topLevel = false

	cv.setJumpPoints(c.Positive)

	switch et := c.Left.(type) {
	case tree.Argument:
		ix := argument[et.Index]
		for i, l := range c.Rights {

			var n label
			if i != len(c.Rights)-1 {
				n = cv.goToNextComparison(false)
			}

			switch k := l.(type) {
			case tree.NumericLiteral:
				cv.jumpOnK(k.Value, ix, tree.EQL)
			case tree.Argument:
				rx := argument[k.Index]
				cv.jumpOnX(ix, rx, tree.EQL)
			}
			if i != len(c.Rights)-1 {
				cv.setJumpPoints(c.Positive)
				cv.c.labelHere(n)
			}
		}
	case tree.NumericLiteral:
		for i, l := range c.Rights {

			var n label
			if i != len(c.Rights)-1 {
				n = cv.goToNextComparison(false)
			}

			k := l.(tree.Argument)
			ix := argument[k.Index]
			cv.jumpOnK(et.Value, ix, tree.EQL)

			if i != len(c.Rights)-1 {
				cv.setJumpPoints(c.Positive)
				cv.c.labelHere(n)
			}
		}
	}
}

func (cv *compilerVisitor) AcceptNegation(c tree.Negation) {
	cv.topLevel = false
	a := &compilerVisitor{c: cv.c, topLevel: false, jf: cv.jt, jt: cv.jf}
	c.Operand.Accept(a)
}

func (cv *compilerVisitor) AcceptNumericLiteral(l tree.NumericLiteral) {
}

func (cv *compilerVisitor) AcceptAnd(c tree.And) {
	n := nextLabel()
	a := &compilerVisitor{c: cv.c, topLevel: false, jf: cv.jf, jt: n}
	c.Left.Accept(a)
	cv.c.labelHere(n)
	c.Right.Accept(cv)
}

func (cv *compilerVisitor) AcceptOr(c tree.Or) {
	n := nextLabel()
	cv.topLevel = false
	a := &compilerVisitor{c: cv.c, topLevel: false, jf: n, jt: cv.jt}
	c.Left.Accept(a)
	cv.c.labelHere(n)
	c.Right.Accept(cv)
}

func (cv *compilerVisitor) AcceptVariable(tree.Variable) {
	panic(fmt.Sprintf("Programming error: there should never be any unexpanded variables if the unifier works correctly: syscall: %s - %s", cv.c.currentlyCompilingSyscall, tree.ExpressionString(cv.c.currentlyCompilingExpression)))
}
