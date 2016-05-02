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

func (cv *compilerVisitor) getLower(k uint64) uint32 {
	return uint32(k)
}

func (cv *compilerVisitor) getUpper(k uint64) uint32 {
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
	cv.c.jumpOnXComparison(op, true, false, cv.negated, false)
	cv.c.loadAt(lx.lower)
	cv.c.jumpOnXComparison(op, true, true, cv.negated, false)
}

func (cv *compilerVisitor) AcceptComparison(c tree.Comparison) {
	cv.topLevel = false
	argL, argR, litL, litR, leftArg, rightArg, leftLit, rightLit := detectSpecialCases(c)

	if leftArg && rightLit {
		ix := argument[argL.Index]
		cv.compareArgToNumeric(litR.Value, ix, c.Op, true)
	}

	if leftLit && rightArg {
		ix := argument[argR.Index]
		cv.compareArgToNumeric(litL.Value, ix, c.Op, true)
	}

	if leftArg && rightArg {
		rx := argument[argR.Index]
		lx := argument[argL.Index]

		cv.jumpOnXChained(rx, lx, c.Op, cv.getChainedJumps(hiTerm), cv.getChainedJumps(lowTerm))
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
		cv.c.jumpOnXComparison(c.Op, cv.terminalJF, cv.terminalJT, cv.negated, false)
	}
}

type jumpType string

const (
	hiTerm  = "hiTerm"
	lowTerm = "lowTerm"
	hi      = "hi"
	low     = "low"
	negHi   = "negHi"
	negLow  = "negLow"
)

// TODO on init
func (cv *compilerVisitor) getChainedJumps(j jumpType) map[jumps]bool {
	hiT := map[jumps]bool{jf: true, jt: false, neg: cv.negated, chained: false}
	lowT := map[jumps]bool{jf: true, jt: true, neg: cv.negated, chained: false}
	hiJ := map[jumps]bool{jf: true, jt: true, neg: cv.negated, chained: true}
	lowJ := map[jumps]bool{jf: false, jt: true, neg: cv.negated, chained: false}
	negH := map[jumps]bool{jf: true, jt: true, neg: cv.negated, chained: true}
	negL := map[jumps]bool{jf: true, jt: false, neg: cv.negated, chained: false}

	allPoints := map[jumpType]map[jumps]bool{
		hiTerm:  hiT,
		lowTerm: lowT,
		hi:      hiJ,
		low:     lowJ,
		negHi:   negH,
		negLow:  negL,
	}
	return allPoints[j]
}

func (cv *compilerVisitor) jumpOnK(l uint64, ix argumentPosition, op tree.ComparisonType, hiJumps map[jumps]bool, lowJumps map[jumps]bool) {
	cv.c.loadAt(ix.upper)
	cv.c.jumpOnKComparison(cv.getUpper(l), op, hiJumps[jf], hiJumps[jt], hiJumps[neg], hiJumps[chained])
	cv.c.loadAt(ix.lower)
	cv.c.jumpOnKComparison(cv.getLower(l), op, lowJumps[jf], lowJumps[jt], lowJumps[neg], lowJumps[chained])
}

func (cv *compilerVisitor) compareArgToNumeric(l uint64, ix argumentPosition, op tree.ComparisonType, isLast bool) {
	if isLast {
		cv.jumpOnK(l, ix, op, cv.getChainedJumps(hiTerm), cv.getChainedJumps(lowTerm))
	} else {
		if cv.negated {
			cv.jumpOnK(l, ix, op, cv.getChainedJumps(negHi), cv.getChainedJumps(negLow))
		} else {
			cv.jumpOnK(l, ix, op, cv.getChainedJumps(hi), cv.getChainedJumps(low))
		}
	}
}

func (cv *compilerVisitor) jumpOnXChained(ix argumentPosition, rx argumentPosition, op tree.ComparisonType, hiJumps map[jumps]bool, lowJumps map[jumps]bool) {
	cv.c.loadAt(ix.upper)
	cv.c.moveAtoX()
	cv.c.loadAt(rx.upper)
	cv.c.jumpOnXComparison(op, hiJumps[jf], hiJumps[jt], hiJumps[neg], hiJumps[chained])

	cv.c.loadAt(ix.lower)
	cv.c.moveAtoX()
	cv.c.loadAt(rx.lower)
	cv.c.jumpOnXComparison(op, lowJumps[jf], lowJumps[jt], lowJumps[neg], lowJumps[chained])
}

func (cv *compilerVisitor) AcceptInclusion(c tree.Inclusion) {
	cv.topLevel = false
	if !c.Positive {
		cv.negated = true
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
					cv.jumpOnXChained(ix, rx, tree.EQL, cv.getChainedJumps(hiTerm), cv.getChainedJumps(lowTerm))
				} else {
					if cv.negated {
						cv.jumpOnXChained(ix, rx, tree.EQL, cv.getChainedJumps(negHi), cv.getChainedJumps(negLow))
					} else {
						cv.jumpOnXChained(ix, rx, tree.EQL, cv.getChainedJumps(hi), cv.getChainedJumps(low))
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
