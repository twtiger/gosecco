package simplify

import "github.com/twtiger/gosecco/tree"

// AcceptNegation implements Visitor
func (s *simplifier) AcceptNegation(v tree.Negation) {
	val := Simplify(v.Operand)
	val2, ok := potentialExtractBooleanValue(val)
	if ok {
		s.result = tree.BooleanLiteral{!val2}
	}
}
