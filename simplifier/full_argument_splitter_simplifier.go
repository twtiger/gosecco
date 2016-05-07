package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptComparison implements Visitor
func (s *fullArgumentSplitterSimplifier) AcceptComparison(a tree.Comparison) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)

	pral, okal := potentialExtractFullArgument(l)
	prnl, oknl := potentialExtractValue(l)

	prar, okar := potentialExtractFullArgument(r)
	prnr, oknr := potentialExtractValue(r)

	if okal && oknr {
		switch a.Op {
		case tree.EQL, tree.NEQL:
			s.result = tree.And{
				Left:  tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.NumericLiteral{prnr & 0xFFFFFFFF}},
				Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.NumericLiteral{(prnr >> 32) & 0xFFFFFFFF}},
			}
		default:
			s.result = tree.Comparison{Op: a.Op, Left: l, Right: r}
		}
	} else if okar && oknl {
		switch a.Op {
		case tree.EQL, tree.NEQL:
			s.result = tree.And{
				Left:  tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{prnl & 0xFFFFFFFF}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				Right: tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{(prnl >> 32) & 0xFFFFFFFF}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
			}
		default:
			s.result = tree.Comparison{Op: a.Op, Left: l, Right: r}
		}

	} else {
		s.result = tree.Comparison{Op: a.Op, Left: l, Right: r}
	}
}

// fullArgumentSplitterSimplifier simplifies full argument references in such a way that
// after this has run, there will be no references to full arguments
// this simplifier is expected to run after the inclusion simplifiers and the LT and LTE simplifiers
// since it will not deal well with those situations
// It can compare full arguments against each other
// It can also deal well with arguments on one side and numbers on the other side
// If the result on one side is the result of a calculation, this simplifier
// will default to assume the wanted behavior is that the upper half of the other side is
// all zeroes. Everything else is obvious.
// It deals specifically with the cases for EQL, NEQL, GT and GTE
type fullArgumentSplitterSimplifier struct {
	nullSimplifier
}

func createFullArgumentSplitterSimplifier() Simplifier {
	s := &fullArgumentSplitterSimplifier{}
	s.realSelf = s
	return s
}
