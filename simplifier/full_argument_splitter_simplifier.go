package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptComparison implements Visitor
func (s *fullArgumentSplitterSimplifier) AcceptComparison(a tree.Comparison) {
	l := s.Simplify(a.Left)
	r := s.Simplify(a.Right)

	pral, okal := potentialExtractFullArgument(l)
	prnlLow, prnlHi, oknl := potentialExtractValueParts(l)

	prar, okar := potentialExtractFullArgument(r)
	prnrLow, prnrHi, oknr := potentialExtractValueParts(r)

	if okal && oknr {
		switch a.Op {
		case tree.EQL:
			s.result = tree.And{
				Left:  tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.NumericLiteral{prnrLow}},
				Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.NumericLiteral{prnrHi}},
			}
		case tree.NEQL:
			s.result = tree.Or{
				Left:  tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.NumericLiteral{prnrLow}},
				Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.NumericLiteral{prnrHi}},
			}
		case tree.GT, tree.GTE:
			s.result = tree.Or{
				Left: tree.Comparison{Op: tree.GT, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.NumericLiteral{prnrHi}},
				Right: tree.And{
					Left:  tree.Comparison{Op: tree.EQL, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.NumericLiteral{prnrHi}},
					Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.NumericLiteral{prnrLow}},
				},
			}
		default:
			panic("shouldn't happen")
		}
	} else if okar && oknl {
		switch a.Op {
		case tree.EQL:
			s.result = tree.And{
				Left:  tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{prnlLow}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				Right: tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{prnlHi}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
			}
		case tree.NEQL:
			s.result = tree.Or{
				Left:  tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{prnlLow}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				Right: tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{prnlHi}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
			}
		case tree.GT, tree.GTE:
			s.result = tree.Or{
				Left: tree.Comparison{Op: tree.GT, Left: tree.NumericLiteral{prnlHi}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
				Right: tree.And{
					Left:  tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{prnlHi}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
					Right: tree.Comparison{Op: a.Op, Left: tree.NumericLiteral{prnlLow}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				},
			}
		default:
			panic("shouldn't happen")
		}
	} else if okal && okar {
		switch a.Op {
		case tree.EQL:
			s.result = tree.And{
				Left:  tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
			}
		case tree.NEQL:
			s.result = tree.Or{
				Left:  tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
			}
		case tree.GT, tree.GTE:
			s.result = tree.Or{
				Left: tree.Comparison{Op: tree.GT, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
				Right: tree.And{
					Left:  tree.Comparison{Op: tree.EQL, Left: tree.Argument{Type: tree.Hi, Index: pral}, Right: tree.Argument{Type: tree.Hi, Index: prar}},
					Right: tree.Comparison{Op: a.Op, Left: tree.Argument{Type: tree.Low, Index: pral}, Right: tree.Argument{Type: tree.Low, Index: prar}},
				},
			}
		default:
			panic("shouldn't happen")
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
