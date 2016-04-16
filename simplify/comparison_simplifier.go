package simplify

import "github.com/twtiger/gosecco/tree"

// AcceptComparison implements Visitor
func (s *simplifier) AcceptComparison(a tree.Comparison) {
	l := Simplify(a.Left)
	r := Simplify(a.Right)

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
		case tree.BIT:
			s.result = tree.BooleanLiteral{(pl & pr) != 0}
			return
		}
	}
	s.result = tree.Comparison{Op: a.Op, Left: l, Right: r}

}
