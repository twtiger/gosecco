package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptAnd implements Visitor
func (s *simplifier) AcceptAnd(a tree.And) {
	l := Simplify(a.Left)
	r := Simplify(a.Right)
	pl, ok1 := potentialExtractBooleanValue(l)
	pr, ok2 := potentialExtractBooleanValue(r)
	// First branch is possible to calculate at compile time
	if ok1 {
		if pl {
			// If the first branch is always true, we are determined by the second branch
			if ok2 {
				s.result = tree.BooleanLiteral{pr}
			} else {
				s.result = r
			}
		} else {
			// If the first branch is always false, we can never succeed
			s.result = tree.BooleanLiteral{false}
		}
	} else {
		// Second branch is possible to calculate at compile time
		if ok2 {
			if pr {
				// If the second branch statically evaluates to true, the and expression is determined by the left arm
				s.result = l
			} else {
				// And if the second branch is false, it doesn't matter what the first branch is
				s.result = tree.BooleanLiteral{false}
			}
		} else {
			s.result = tree.And{l, r}
		}
	}
}
