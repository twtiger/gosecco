package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptBinaryNegation implements Visitor
func (s *binaryNegationSimplifier) AcceptBinaryNegation(v tree.BinaryNegation) {
	val := s.Simplify(v.Operand)
	s.result = tree.Arithmetic{Op: tree.BINXOR, Left: val, Right: tree.NumericLiteral{uint64(0xFFFFFFFFFFFFFFFF)}}
}

// binaryNegationSimplifier simplifies binary complement by removing it and replacing it with an xor instruction
type binaryNegationSimplifier struct {
	nullSimplifier
}

func createBinaryNegationSimplifier() Simplifier {
	s := &binaryNegationSimplifier{}
	s.realSelf = s
	return s
}
