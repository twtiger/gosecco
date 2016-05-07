package simplifier

import "github.com/twtiger/gosecco/tree"

// Simplify will take an expression and reduce it as much as possible using state operations
func Simplify(inp tree.Expression) tree.Expression {
	s := &simplifier{inp}
	inp.Accept(s)
	return s.result
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
