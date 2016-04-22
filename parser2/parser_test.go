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

func (s *ParserSuite) Test_parsesDivision(c *C) {
	result := parseExpression("42 / 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.DIV, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesModulo(c *C) {
	result := parseExpression("42 % 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.MOD, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesRSH(c *C) {
	result := parseExpression("42 >> 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.RSH, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesLSH(c *C) {
	result := parseExpression("42 << 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.LSH, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesOR(c *C) {
	result := parseExpression("42 | 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.BINOR, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesAND(c *C) {
	result := parseExpression("42 & 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.BINAND, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesXOR(c *C) {
	result := parseExpression("42 ^ 15")

	c.Assert(result, DeepEquals, tree.Arithmetic{Op: tree.BINXOR, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}})
}

func (s *ParserSuite) Test_parsesMultiplicationAndAddition(c *C) {
	result := parseExpression("42 * 15 + 1")

	c.Assert(result, DeepEquals,
		tree.Arithmetic{Op: tree.PLUS, Left: tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{15}}, Right: tree.NumericLiteral{1}})
}

func (s *ParserSuite) Test_parsesParens(c *C) {
	result := parseExpression("42 * (15 + 1)")

	c.Assert(result, DeepEquals,
		tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{15}, Right: tree.NumericLiteral{1}}})
}

func (s *ParserSuite) Test_parsesArgument(c *C) {
	result := parseExpression("42 * arg1")

	c.Assert(result, DeepEquals,
		tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.Argument{1}})
}

func (s *ParserSuite) Test_parsesVariable(c *C) {
	result := parseExpression("42 * arg6")

	c.Assert(result, DeepEquals,
		tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.Variable{"arg6"}})
}

func (s *ParserSuite) Test_parsesUnaryNegation(c *C) {
	result := parseExpression("42 * ~arg6")

	c.Assert(result, DeepEquals,
		tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{42}, Right: tree.BinaryNegation{tree.Variable{"arg6"}}})
}

func (s *ParserSuite) Test_parsesBooleanExpression(c *C) {
	result := parseExpression("true && false")
	c.Assert(result, DeepEquals, tree.And{Left: tree.BooleanLiteral{true}, Right: tree.BooleanLiteral{false}})

	result = parseExpression("true || false")
	c.Assert(result, DeepEquals, tree.Or{Left: tree.BooleanLiteral{true}, Right: tree.BooleanLiteral{false}})

	result = parseExpression("!(true || false)")
	c.Assert(result, DeepEquals, tree.Negation{tree.Or{Left: tree.BooleanLiteral{true}, Right: tree.BooleanLiteral{false}}})
}

func (s *ParserSuite) Test_parsesEquality(c *C) {
	result := parseExpression("true == false")
	c.Assert(result, DeepEquals, tree.Comparison{Op: tree.EQL, Left: tree.BooleanLiteral{true}, Right: tree.BooleanLiteral{false}})

	result = parseExpression("42 != 1")
	c.Assert(result, DeepEquals, tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}})

	result = parseExpression("42 > 1")
	c.Assert(result, DeepEquals, tree.Comparison{Op: tree.GT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}})

	result = parseExpression("42 >= 1")
	c.Assert(result, DeepEquals, tree.Comparison{Op: tree.GTE, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}})

	result = parseExpression("42 < 1")
	c.Assert(result, DeepEquals, tree.Comparison{Op: tree.LT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}})

	result = parseExpression("42 <= 1")
	c.Assert(result, DeepEquals, tree.Comparison{Op: tree.LTE, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}})
}

func (s *ParserSuite) Test_parseCall(c *C) {
	result := parseExpression("foo(1+1, 2+3, arg0)")
	c.Assert(result, DeepEquals,
		tree.Call{Name: "foo",
			Args: []tree.Any{
				tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x1}, Right: tree.NumericLiteral{Value: 0x1}},
				tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x2}, Right: tree.NumericLiteral{Value: 0x3}},
				tree.Argument{Index: 0}}})
}

func (s *ParserSuite) Test_parseIn(c *C) {
	result := parseExpression("in(1+1, 2+3, arg0)")
	c.Assert(result, DeepEquals,
		tree.Inclusion{Positive: true,
			Left: tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x1}, Right: tree.NumericLiteral{Value: 0x1}},
			Rights: []tree.Numeric{
				tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x2}, Right: tree.NumericLiteral{Value: 0x3}},
				tree.Argument{Index: 0}}})
}

func (s *ParserSuite) Test_parseNotIn(c *C) {
	result := parseExpression("notin(1+1, 2+3, arg0)")
	c.Assert(result, DeepEquals,
		tree.Inclusion{Positive: false,
			Left: tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x1}, Right: tree.NumericLiteral{Value: 0x1}},
			Rights: []tree.Numeric{
				tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x2}, Right: tree.NumericLiteral{Value: 0x3}},
				tree.Argument{Index: 0}}})
}
