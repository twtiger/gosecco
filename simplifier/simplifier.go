package simplifier

import "github.com/twtiger/gosecco/tree"

// Simplifier is something that can simplify expression
type Simplifier interface {
	tree.Visitor
	Simplify(tree.Expression) tree.Expression
}

func reduceSimplify(inp tree.Expression, ss ...Simplifier) tree.Expression {
	result := inp

	for _, s := range ss {
		result = s.Simplify(result)
	}

	return result
}

// Simplify will take an expression and reduce it as much as possible using state operations
func Simplify(inp tree.Expression) tree.Expression {
	return reduceSimplify(inp,
		createLtExpressionsSimplifier(),
		createArithmeticSimplifier(),
		createComparisonSimplifier(),
		createBooleanSimplifier(),
		createBinaryNegationSimplifier(),
		createInclusionSimplifier(),
	)
}

type simplifier struct {
	result tree.Expression
}

func potentialExtractValue(a tree.Numeric) (uint64, bool) {
	v, ok := a.(tree.NumericLiteral)
	if ok {
		return v.Value, ok
	}
	return 0, false
}

func potentialExtractBooleanValue(a tree.Boolean) (bool, bool) {
	v, ok := a.(tree.BooleanLiteral)
	if ok {
		return v.Value, ok
	}
	return false, false
}
