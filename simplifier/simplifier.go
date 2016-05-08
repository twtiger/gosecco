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
		// X in [P]  ==>  P == Q
		// X in [P, Q, R]  where X and R can be determined to not be equal  ==>  X in [P, Q]
		// X in [P, Q, R]  where X and one of the values can be determined to be equal  ==>  true
		// X notIn [P]  ==>  P != Q
		// X notIn [P, Q, R]  where X and R can be determined to not be equal  ==>  X notIn [P, Q]
		// X notIn [P, Q, R]  where X and one of the values can be determined to be equal  ==>  false
		createInclusionSimplifier(),

		// X in [P, Q, R]     ==>  X == P || X == Q || X == R
		// X notIn [P, Q, R]  ==>  X != P && X != Q && X != R
		createInclusionRemoverSimplifier(),

		// X < Y    ==>  Y >= X
		// X <= Y   ==>  Y > X
		createLtExpressionsSimplifier(),

		// Where X and Y can be determined statically:
		// X + Y   ==>  [X+Y]
		// X - Y   ==>  [X-Y]
		// X * Y   ==>  [X*Y]
		// X / Y   ==>  [X/Y]
		// X % Y   ==>  [X%Y]
		// X & Y   ==>  [X&Y]
		// X | Y   ==>  [X|Y]
		// X ^ Y   ==>  [X^Y]
		// X << Y  ==>  [X<<Y]
		// X >> Y  ==>  [X<<Y]
		// ~X      ==>  [~X]
		// Note that these calculations will all be done on 64bit unsigned values
		// - this could lead to different result than if they were evaluated by the BPF engine.
		createArithmeticSimplifier(),

		// Where X and Y can be determined statically:
		// X == Y  where X == Y  ==>  true
		// X == Y  where X != Y  ==>  false
		// X != Y  where X == Y  ==>  false
		// X != Y  where X != Y  ==>  true
		// X > Y   where X > Y   ==>  true
		// X > Y   where X <= Y  ==>  false
		// X >= Y  where X >= Y  ==>  true
		// X >= Y  where X < Y   ==>  false
		// X < Y   where X < Y   ==>  true
		// X < Y   where X >= Y  ==>  false
		// X <= Y  where X <= Y  ==>  true
		// X <= Y  where X > Y   ==>  false
		createComparisonSimplifier(),

		// !true           ==>  false
		// !false          ==>  true
		// false || Y      ==>  Y
		// false || true   ==>  true
		// false || false  ==>  false
		// true  || Y      ==>  true
		// true  && true   ==>  true
		// true  && false  ==>  false
		// true  && Y      ==>  Y
		// false && [any]  ==>  false
		createBooleanSimplifier(),

		// ~X  ==> X ^ 0xFFFFFFFFFFFFFFFF
		createBinaryNegationSimplifier(),

		// Where X can be determined statically (the opposite order is also valid, and X can also be an arg)
		// arg0 == X  ==>  argL0 == X.low && argH0 == X.high
		// arg0 != X  ==>  argL0 != X.low || argH0 != X.high
		// arg0 > X   ==>  argH0 > X.high || (argH0 == X.high && argL0 > X.low)
		// arg0 >= X  ==>  argH0 > X.high || (argH0 == X.high && argL0 >= X.low)
		createFullArgumentSplitterSimplifier(),
	)
}

type simplifier struct {
	result tree.Expression
}

func potentialExtractFullArgument(a tree.Expression) (int, bool) {
	v, ok := a.(tree.Argument)
	if ok && v.Type == tree.Full {
		return v.Index, ok
	}
	return 0, false
}

func potentialExtractValue(a tree.Numeric) (uint64, bool) {
	v, ok := a.(tree.NumericLiteral)
	if ok {
		return v.Value, ok
	}
	return 0, false
}

func potentialExtractValueParts(a tree.Numeric) (uint64, uint64, bool) {
	v, ok := a.(tree.NumericLiteral)
	if ok {
		low := v.Value & 0xFFFFFFFF
		high := (v.Value >> 32) & 0xFFFFFFFF
		return low, high, ok
	}
	return 0, 0, false
}

func potentialExtractBooleanValue(a tree.Boolean) (bool, bool) {
	v, ok := a.(tree.BooleanLiteral)
	if ok {
		return v.Value, ok
	}
	return false, false
}
