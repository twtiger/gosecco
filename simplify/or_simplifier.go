package simplify

import "github.com/twtiger/gosecco/tree"

// AcceptOr implements Visitor
func (s *simplifier) AcceptOr(a tree.Or) {
	l := Simplify(a.Left)
	r := Simplify(a.Right)
	pl, ok1 := potentialExtractBooleanValue(l)
	pr, ok2 := potentialExtractBooleanValue(r)
	// First branch is possible to calculate at compile time
	if ok1 {
		if pl {
			s.result = tree.BooleanLiteral{true}
		} else {
			if ok2 {
				s.result = tree.BooleanLiteral{pr}
			} else {
				s.result = r
			}
		}
	} else {
		s.result = tree.Or{l, r}
	}
}
