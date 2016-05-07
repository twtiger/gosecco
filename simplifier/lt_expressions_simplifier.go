package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptComparison implements Visitor
func (s *ltExpressionsSimplifier) AcceptComparison(a tree.Comparison) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)

	newOp := a.Op

	switch a.Op {
	case tree.LT:
		newOp = tree.GTE
		l, r = r, l
	case tree.LTE:
		newOp = tree.GT
		l, r = r, l
	}

	s.result = tree.Comparison{Op: newOp, Left: l, Right: r}

}

// ltExpressionsSimplifier simplifies LT and LTE expressions by rewriting them to GT and GTE expressions
type ltExpressionsSimplifier struct {
	nullSimplifier
}

func createLtExpressionsSimplifier() Simplifier {
	s := &ltExpressionsSimplifier{}
	s.realSelf = s
	return s
}
