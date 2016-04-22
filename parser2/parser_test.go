package parser2

import (
	"testing"

	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ParserSuite struct{}

var _ = Suite(&ParserSuite{})

func (s *ParserSuite) Test_parsesNumber(c *C) {
	result := parseExpression("42")

	c.Assert(result, DeepEquals, tree.NumericLiteral{42})
}

func (s *ParserSuite) Test_parsesAddition(c *C) {
	result := parseExpression("42 + 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesMultiplication(c *C) {
	result := parseExpression("42 * 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesMultiplicationAndAddition(c *C) {
	result := parseExpression("42 * 15 + 1")

	c.Assert(result, DeepEquals,
		tree.Arithmetic{Op: tree.PLUS, Left: tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}}, Right: tree.NumericLiteral{1}})
}
