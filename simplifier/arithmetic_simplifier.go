package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptBinaryNegation implements Visitor
func (s *arithmeticSimplifier) AcceptBinaryNegation(v tree.BinaryNegation) {
	val := s.Simplify(v.Operand)
	val2, ok := potentialExtractValue(val)
	if ok {

		s.result = tree.NumericLiteral{^val2}
	}
}

// AcceptArithmetic implements Visitor
func (s *arithmeticSimplifier) AcceptArithmetic(a tree.Arithmetic) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)

	pl, ok1 := potentialExtractValue(l)
	pr, ok2 := potentialExtractValue(r)

	if ok1 && ok2 {
		switch a.Op {
		case tree.PLUS:
			s.result = tree.NumericLiteral{pl + pr}
			return
		case tree.MINUS:
			s.result = tree.NumericLiteral{pl - pr}
			return
		case tree.MULT:
			s.result = tree.NumericLiteral{pl * pr}
			return
		case tree.DIV:
			s.result = tree.NumericLiteral{pl / pr}
			return
		case tree.MOD:
			s.result = tree.NumericLiteral{pl % pr}
			return
		case tree.BINAND:
			s.result = tree.NumericLiteral{pl & pr}
			return
		case tree.BINOR:
			s.result = tree.NumericLiteral{pl | pr}
			return
		case tree.BINXOR:
			s.result = tree.NumericLiteral{pl ^ pr}
			return
		case tree.LSH:
			s.result = tree.NumericLiteral{pl << pr}
			return
		case tree.RSH:
			s.result = tree.NumericLiteral{pl >> pr}
			return
		}
	}
	s.result = tree.Arithmetic{Op: a.Op, Left: l, Right: r}
}

// arithmeticSimplifier simplifies arithmetic expressions by calculating them as much as possible
type arithmeticSimplifier struct {
	nullSimplifier
}

func createArithmeticSimplifier() Simplifier {
	s := &arithmeticSimplifier{}
	s.realSelf = s
	return s
}
