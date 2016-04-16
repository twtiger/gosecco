package simplify

import "github.com/twtiger/gosecco/tree"

// AcceptBinaryNegation implements Visitor
func (s *simplifier) AcceptBinaryNegation(v tree.BinaryNegation) {
	val := Simplify(v.Operand)
	val2, ok := potentialExtractValue(val)
	if ok {

		s.result = tree.NumericLiteral{^val2}
	}
}
