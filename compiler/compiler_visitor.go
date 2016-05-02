package compiler

import (
	"fmt"

	"github.com/twtiger/gosecco/tree"
)

type compilerVisitor struct {
	c         *compiler
	terminal  bool
	exclusive bool
	negated   bool
	inverted  bool
	topLevel  bool
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
	cv.c.jumpOnXComparison(op, true, false, cv.inverted, false)
	cv.c.loadAt(lx.lower)
	cv.c.jumpOnXComparison(op, true, true, cv.inverted, false)
}

func (cv *compilerVisitor) AcceptComparison(c tree.Comparison) {
	cv.topLevel = false
	argL, argR, litL, litR, leftArg, rightArg, leftLit, rightLit := detectSpecialCases(c)

	if leftArg && rightLit {
		ix := argument[argL.Index]
		cv.compareArgToNumeric(litR.Value, ix, c.Op, cv.terminal)
	}

	if leftLit && rightArg {
		ix := argument[argR.Index]
		cv.compareArgToNumeric(litL.Value, ix, c.Op, cv.terminal)
	}

	if leftArg && rightArg {
		rx := argument[argR.Index]
		lx := argument[argL.Index]

		cv.jumpOnXChained(rx, lx, c.Op, jumpPoints[hiTerm], jumpPoints[lowTerm])
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
		cv.c.jumpOnXComparison(c.Op, true, true, cv.inverted, false)
	}
}

type jumpType string

const (
	hiTerm  = "hiTerm"
	lowTerm = "lowTerm"
	hi      = "hi"
	low     = "low"
	negHi   = "hiNeg"
	negLow  = "lowNeg"
	invHi   = "negHi"
	invLow  = "negLow"
	hiExl   = "hiExl"
	lowExl  = "lowExl"
)

var jumpPoints = map[jumpType]map[jumps]bool{
	hiTerm:  map[jumps]bool{jf: true, jt: false, chained: false},
	lowTerm: map[jumps]bool{jf: true, jt: true, chained: false},
	hi:      map[jumps]bool{jf: true, jt: true, chained: true},
	low:     map[jumps]bool{jf: false, jt: true, chained: false},
	hiExl:   map[jumps]bool{jf: true, jt: true, chained: true},
	lowExl:  map[jumps]bool{jf: true, jt: false, chained: false},
	invHi:   map[jumps]bool{jf: true, jt: true, chained: true},
	invLow:  map[jumps]bool{jf: true, jt: false, chained: false},
	negHi:   map[jumps]bool{jf: true, jt: false, chained: false},
	negLow:  map[jumps]bool{jf: true, jt: false, chained: false},
}

func (cv *compilerVisitor) jumpOnK(l uint64, ix argumentPosition, op tree.ComparisonType, hiJumps map[jumps]bool, lowJumps map[jumps]bool) {
	cv.c.loadAt(ix.upper)
	cv.c.jumpOnKComparison(getUpper(l), op, hiJumps[jf], hiJumps[jt], cv.negated, cv.inverted, hiJumps[chained])
	cv.c.loadAt(ix.lower)
	cv.c.jumpOnKComparison(getLower(l), op, lowJumps[jf], lowJumps[jt], cv.negated, cv.inverted, lowJumps[chained])
}

func (cv *compilerVisitor) compareArgToNumeric(l uint64, ix argumentPosition, op tree.ComparisonType, isLast bool) {
	switch {
	case isLast:
		cv.jumpOnK(l, ix, op, jumpPoints[hiTerm], jumpPoints[lowTerm])
	case !cv.exclusive && cv.negated:
		cv.jumpOnK(l, ix, op, jumpPoints[negHi], jumpPoints[negLow])
	case cv.inverted:
		cv.jumpOnK(l, ix, op, jumpPoints[invHi], jumpPoints[invLow])
	case cv.exclusive && !cv.terminal && !cv.negated:
		cv.jumpOnK(l, ix, op, jumpPoints[hiExl], jumpPoints[lowExl])
	default:
		cv.jumpOnK(l, ix, op, jumpPoints[hi], jumpPoints[low])
	}
}

func (cv *compilerVisitor) jumpOnXChained(ix argumentPosition, rx argumentPosition, op tree.ComparisonType, hiJumps map[jumps]bool, lowJumps map[jumps]bool) {
	cv.c.loadAt(ix.upper)
	cv.c.moveAtoX()
	cv.c.loadAt(rx.upper)
	cv.c.jumpOnXComparison(op, hiJumps[jf], hiJumps[jt], cv.inverted, hiJumps[chained])

	cv.c.loadAt(ix.lower)
	cv.c.moveAtoX()
	cv.c.loadAt(rx.lower)
	cv.c.jumpOnXComparison(op, lowJumps[jf], lowJumps[jt], cv.inverted, lowJumps[chained])
}

func (cv *compilerVisitor) AcceptInclusion(c tree.Inclusion) {
	cv.topLevel = false
	if !c.Positive {
		cv.inverted = true
	}

	switch et := c.Left.(type) {
	case tree.Argument:
		ix := argument[et.Index]
		for i, l := range c.Rights {
			isLast := i == len(c.Rights)-1
			switch k := l.(type) {
			case tree.NumericLiteral:
				cv.compareArgToNumeric(k.Value, ix, tree.EQL, isLast)
			case tree.Argument:
				rx := argument[k.Index]
				if isLast {
					cv.jumpOnXChained(ix, rx, tree.EQL, jumpPoints[hiTerm], jumpPoints[lowTerm])
				} else {
					if cv.negated {
						cv.jumpOnXChained(ix, rx, tree.EQL, jumpPoints[negHi], jumpPoints[negLow])
					} else {
						cv.jumpOnXChained(ix, rx, tree.EQL, jumpPoints[hi], jumpPoints[low])
					}
				}
			}
		}
	case tree.NumericLiteral:
		for i, l := range c.Rights {
			k := l.(tree.Argument)
			ix := argument[k.Index]
			isLast := i == len(c.Rights)-1
			cv.compareArgToNumeric(et.Value, ix, tree.EQL, isLast)
		}
	}
}

func (cv *compilerVisitor) AcceptNegation(c tree.Negation) {
	cv.topLevel = false
	cv.negated = true
	c.Operand.Accept(cv)
}

func (cv *compilerVisitor) AcceptNumericLiteral(l tree.NumericLiteral) {
}

func (cv *compilerVisitor) AcceptAnd(c tree.And) {
	cv.topLevel = false
	cv.exclusive = true
	cv.terminal = false
	c.Left.Accept(cv)
	cv.terminal = true
	c.Right.Accept(cv)
}

func (cv *compilerVisitor) AcceptOr(c tree.Or) {
	cv.topLevel = false
	cv.terminal = false
	c.Left.Accept(cv)
	cv.terminal = true
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
