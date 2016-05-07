package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptComparison implements Visitor
func (s *comparisonSimplifier) AcceptComparison(a tree.Comparison) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)

	pl, ok1 := potentialExtractValue(l)
	pr, ok2 := potentialExtractValue(r)

	if ok1 && ok2 {
		switch a.Op {
		case tree.EQL:
			s.result = tree.BooleanLiteral{pl == pr}
			return
		case tree.NEQL:
			s.result = tree.BooleanLiteral{pl != pr}
			return
		case tree.GT:
			s.result = tree.BooleanLiteral{pl > pr}
			return
		case tree.GTE:
			s.result = tree.BooleanLiteral{pl >= pr}
			return
		case tree.LT:
			s.result = tree.BooleanLiteral{pl < pr}
			return
		case tree.LTE:
			s.result = tree.BooleanLiteral{pl <= pr}
			return
		}
	}
	s.result = tree.Comparison{Op: a.Op, Left: l, Right: r}
}

// comparisonSimplifier simplifies comparison expressions by calculating them as much as possible
type comparisonSimplifier struct {
	nullSimplifier
}

func createComparisonSimplifier() Simplifier {
	s := &comparisonSimplifier{}
	s.realSelf = s
	return s
}
