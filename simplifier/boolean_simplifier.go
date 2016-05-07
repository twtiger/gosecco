package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptOr implements Visitor
func (s *booleanSimplifier) AcceptOr(a tree.Or) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)
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

// AcceptAnd implements Visitor
func (s *booleanSimplifier) AcceptAnd(a tree.And) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)
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

// AcceptNegation implements Visitor
func (s *booleanSimplifier) AcceptNegation(v tree.Negation) {
	val := s.Simplify(v.Operand)
	val2, ok := potentialExtractBooleanValue(val)
	if ok {
		s.result = tree.BooleanLiteral{!val2}
	}
}

// booleanSimplifier simplifies boolean expressions by calculating them as much as possible
type booleanSimplifier struct {
	nullSimplifier
}

func createBooleanSimplifier() Simplifier {
	s := &booleanSimplifier{}
	s.realSelf = s
	return s
}
