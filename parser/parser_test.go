package parser

import (
	"github.com/twtiger/go-seccomp/tree"
	. "gopkg.in/check.v1"
)

type RulesSuite struct{}

var _ = Suite(&RulesSuite{})

func (s *RulesSuite) Test_parsesSimpleRule(c *C) {
	result, _ := parseExpression("1")

	c.Assert(result, DeepEquals, tree.BooleanLiteral{true})
}

func (s *RulesSuite) Test_parsesAlmostSimpleRule(c *C) {
	result, _ := parseExpression("arg0 > 0")

	c.Assert(result, DeepEquals, tree.Comparison{
		Left:  tree.Argument{0},
		Op:    tree.GT,
		Right: tree.NumericLiteral{0},
	})
}

func (s *RulesSuite) Test_parseAnotherRule(c *C) {
	result, _ := parseExpression("arg0 == 4")

	c.Assert(result, DeepEquals, tree.Comparison{
		Left:  tree.Argument{0},
		Op:    tree.EQL,
		Right: tree.NumericLiteral{4},
	})
}

func (s *RulesSuite) Test_parseYetAnotherRule(c *C) {
	result, _ := parseExpression("arg0 == 4 || arg0 == 5")

	c.Assert(tree.ExpressionString(result), Equals, "(or (eq arg0 4) (eq arg0 5))")
	c.Assert(result, DeepEquals, tree.Or{
		Left: tree.Comparison{
			Left:  tree.Argument{0},
			Op:    tree.EQL,
			Right: tree.NumericLiteral{4},
		},
		Right: tree.Comparison{
			Left:  tree.Argument{0},
			Op:    tree.EQL,
			Right: tree.NumericLiteral{5},
		},
	})
}

func parseExpectSuccess(c *C, str string) string {
	result, err := parseExpression(str)
	c.Assert(err, IsNil)
	return tree.ExpressionString(result)
}

func (s *RulesSuite) Test_parseExpressionWithMultiplication(c *C) {
	c.Assert(parseExpectSuccess(c, "arg0 == 12 * 3"), Equals, "(eq arg0 (mul 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithAddition(c *C) {
	result, _ := parseExpression("arg0 == 12 + 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (plus 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithDivision(c *C) {
	result, _ := parseExpression("arg0 == 12 / 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (div 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithSubtraction(c *C) {
	result, _ := parseExpression("arg0 == 12 - 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (minus 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryAnd(c *C) {
	result, _ := parseExpression("arg0 == 0 & 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (binand 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryOr(c *C) {
	result, _ := parseExpression("arg0 == 0 | 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (binor 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryXor(c *C) {
	result, _ := parseExpression("arg0 == 0 ^ 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (binxor 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryNegation(c *C) {
	c.Assert(parseExpectSuccess(c, "arg0 == ^0"), Equals, "(eq arg0 (binNeg 0))")
}

func (s *RulesSuite) Test_parseAExpressionLeftShift(c *C) {
	result, _ := parseExpression("arg0 == 2 << 1")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (lsh 2 1))")
}

func (s *RulesSuite) Test_parseAExpressionRightShift(c *C) {
	result, _ := parseExpression("arg0 == (2 >> 1)")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (rsh 2 1))")
}

func (s *RulesSuite) Test_parseAExpressionWithModulo(c *C) {
	result, _ := parseExpression("arg0 == 12 % 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (mod 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithBooleanAnd(c *C) {
	result, _ := parseExpression("arg0 == 0 && arg1 == 0")
	c.Assert(tree.ExpressionString(result), Equals, "(and (eq arg0 0) (eq arg1 0))")
}

func (s *RulesSuite) Test_parseAExpressionWithBooleanNegation(c *C) {
	c.Assert(parseExpectSuccess(c, "!(arg0 == 1)"), Equals, "(not (eq arg0 1))")
}

func (s *RulesSuite) Test_parseAExpressionWithNotEqual(c *C) {
	result, _ := parseExpression("arg0 != 1")
	c.Assert(tree.ExpressionString(result), Equals, "(neq arg0 1)")
}

func (s *RulesSuite) Test_parseAExpressionWithGreaterThanOrEqualTo(c *C) {
	result, _ := parseExpression("arg0 >= 1")
	c.Assert(tree.ExpressionString(result), Equals, "(gte arg0 1)")
}

func (s *RulesSuite) Test_parseAExpressionWithLessThan(c *C) {
	result, _ := parseExpression("arg0 < arg1")
	c.Assert(tree.ExpressionString(result), Equals, "(lt arg0 arg1)")
}

func (s *RulesSuite) Test_parseAExpressionWithLessThanOrEqualTo(c *C) {
	result, _ := parseExpression("arg0 <= arg1")
	c.Assert(tree.ExpressionString(result), Equals, "(lte arg0 arg1)")
}

func (s *RulesSuite) Test_parseAExpressionWithBitSets(c *C) {
	result, _ := parseExpression("arg0 & val")
	c.Assert(tree.ExpressionString(result), Equals, "(bitSet arg0 val)")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusion(c *C) {
	result, _ := parseExpression("in(arg0, 1, 2)")
	c.Assert(tree.ExpressionString(result), Equals, "(in arg0 1 2)")
}

func (s *RulesSuite) Test_parseAExpressionWithExclusion(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("notIn(arg0, 1, 2)")
	c.Assert(tree.ExpressionString(result), Equals, "(notIn arg0 1, 2)")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusionLargerSet(c *C) {
	result, _ := parseExpression("in(arg0, 1, 2, 3, 4)")
	c.Assert(tree.ExpressionString(result), Equals, "(in arg0 1 2 3 4)")
}

func (s *RulesSuite) Test_parseAExpressionWithExclusionLargerSet(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("notin(arg0, 1, 2, 3, 4)")
	c.Assert(tree.ExpressionString(result), Equals, "(notin arg0 1 2 3 4(")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusionWithWhitespace(c *C) {
	result, _ := parseExpression("in(arg0, 1,   2,   3,   4)")
	c.Assert(tree.ExpressionString(result), Equals, "(in arg0 1 2 3 4)")
}

func (s *RulesSuite) Test_parseAExpressionWithTrue(c *C) {
	result, _ := parseExpression("true")
	c.Assert(tree.ExpressionString(result), Equals, "true")
}

func (s *RulesSuite) Test_parseAExpressionWithFalse(c *C) {
	result, _ := parseExpression("false")
	c.Assert(tree.ExpressionString(result), Equals, "false")
}

func (s *RulesSuite) Test_parseAExpressionWith0AsFalse(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("0")
	c.Assert(tree.ExpressionString(result), Equals, "false")
}

func (s *RulesSuite) Test_parseAExpressionWithParens(c *C) {
	result, _ := parseExpression("arg0 == (12 + 3) * 2")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (mul (plus 12 3) 2))")
}

func (s *RulesSuite) Test_parseAExpressionWithNestedOperators(c *C) {
	result, _ := parseExpression("arg0 == 12 + 3 * 2")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (plus 12 (mul 3 2)))")
}

func (s *RulesSuite) Test_parseAExpressionWithInvalidArithmeticOperator(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 == 12 _ 3")
	c.Assert(tree.ExpressionString(result), Equals, "(eq arg0 (add 12 3))")
}

func (s *RulesSuite) Test_parseArgumentsCorrectly_andIncorrectly(c *C) {
	c.Assert(parseExpectSuccess(c, "arg0 == 0"), Equals, "(eq arg0 0)")
	c.Assert(parseExpectSuccess(c, "arg5 == 0"), Equals, "(eq arg5 0)")

	result, _ := parseExpression("arg6 == 0")
	c.Assert(result, DeepEquals, tree.Comparison{
		Left:  tree.Variable{"arg6"},
		Op:    tree.EQL,
		Right: tree.NumericLiteral{0},
	})
}
