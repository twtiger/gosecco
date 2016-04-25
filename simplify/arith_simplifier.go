package simplify

import "github.com/twtiger/gosecco/tree"

// AcceptArithmetic implements Visitor
func (s *simplifier) AcceptArithmetic(a tree.Arithmetic) {
	l := Simplify(a.Left)
	r := Simplify(a.Right)

	pl, ok1 := potentialExtractValue(l)
	pr, ok2 := potentialExtractValue(r)

	if ok1 && ok2 {
		switch a.Op {
		case tree.PLUS:
			s.result = tree.NumericLiteral{Value: pl + pr}
			return
		case tree.MINUS:
			s.result = tree.NumericLiteral{Value: pl - pr}
			return
		case tree.MULT:
			s.result = tree.NumericLiteral{Value: pl * pr}
			return
		case tree.DIV:
			s.result = tree.NumericLiteral{Value: pl / pr}
			return
		case tree.MOD:
			s.result = tree.NumericLiteral{Value: pl % pr}
			return
		case tree.BINAND:
			s.result = tree.NumericLiteral{Value: pl & pr}
			return
		case tree.BINOR:
			s.result = tree.NumericLiteral{Value: pl | pr}
			return
		case tree.BINXOR:
			s.result = tree.NumericLiteral{Value: pl ^ pr}
			return
		case tree.LSH:
			s.result = tree.NumericLiteral{Value: pl << pr}
			return
		case tree.RSH:
			s.result = tree.NumericLiteral{Value: pl >> pr}
			return
		}
	}
	s.result = tree.Arithmetic{Op: a.Op, Left: l, Right: r}
}
