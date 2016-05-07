package simplifier

import (
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

type FullArgumentSplitterSimplifierSuite struct{}

var _ = Suite(&FullArgumentSplitterSimplifierSuite{})

func (s *FullArgumentSplitterSimplifierSuite) Test_simplifiesEqualityWithArgAgainstNumber(c *C) {
	sx := createFullArgumentSplitterSimplifier().Simplify(
		tree.Comparison{
			Op:    tree.EQL,
			Left:  tree.Argument{Type: tree.Full, Index: 2},
			Right: tree.NumericLiteral{0x123456789ABCDEF0},
		},
	)

	c.Assert(tree.ExpressionString(sx), Equals, "(and (eq argL2 2596069104) (eq argH2 305419896))")

	sx = createFullArgumentSplitterSimplifier().Simplify(
		tree.Comparison{
			Op:    tree.EQL,
			Left:  tree.NumericLiteral{0x123456789ABCDEF0},
			Right: tree.Argument{Type: tree.Full, Index: 2},
		},
	)

	c.Assert(tree.ExpressionString(sx), Equals, "(and (eq 2596069104 argL2) (eq 305419896 argH2))")
}

func (s *FullArgumentSplitterSimplifierSuite) Test_simplifiesNonequalityWithArgAgainstNumber(c *C) {
	sx := createFullArgumentSplitterSimplifier().Simplify(
		tree.Comparison{
			Op:    tree.NEQL,
			Left:  tree.Argument{Type: tree.Full, Index: 2},
			Right: tree.NumericLiteral{0x123456789ABCDEF0},
		},
	)

	c.Assert(tree.ExpressionString(sx), Equals, "(and (neq argL2 2596069104) (neq argH2 305419896))")

	sx = createFullArgumentSplitterSimplifier().Simplify(
		tree.Comparison{
			Op:    tree.NEQL,
			Left:  tree.NumericLiteral{0x123456789ABCDEF0},
			Right: tree.Argument{Type: tree.Full, Index: 2},
		},
	)

	c.Assert(tree.ExpressionString(sx), Equals, "(and (neq 2596069104 argL2) (neq 305419896 argH2))")
}
