package parser

import (
	. "gopkg.in/check.v1"
)

type RulesSuite struct{}

var _ = Suite(&RulesSuite{})

func (s *RulesSuite) Test_parsesSimpleRule(c *C) {
	result, _ := parseRule("read: 1")

	c.Assert(result, DeepEquals, rule{
		syscall:    "read",
		expression: trueLiteral{},
	},
	)
}

func (s *RulesSuite) Test_parsesAlmostSimpleRule(c *C) {
	result, _ := parseRule("read2: arg0 > 0")

	c.Assert(result, DeepEquals, rule{
		syscall: "read2",
		expression: comparison{
			left:  argumentNode{index: 0},
			cmp:   "gt",
			right: literalNode{value: 0},
		},
	})
}

func (s *RulesSuite) Test_parseAnotherRule(c *C) {
	result, _ := parseRule("read3: arg0 == 4")

	c.Assert(result, DeepEquals, rule{
		syscall: "read3",
		expression: comparison{
			left:  argumentNode{index: 0},
			cmp:   "eq",
			right: literalNode{value: 4},
		},
	})
}

func (s *RulesSuite) Test_parseYetAnotherRule(c *C) {
	result, _ := parseRule("read4: arg0 == 4 || arg0 == 5")
	//fmt.Printf("%#v\n", result)

	c.Assert(result.expression.String(), Equals, "(lor (eq arg0 4) (eq arg0 5))")
	c.Assert(result, DeepEquals, rule{
		syscall: "read4",
		expression: orExpr{
			left: comparison{
				left:  argumentNode{index: 0},
				cmp:   "eq",
				right: literalNode{value: 4},
			},
			right: comparison{
				left:  argumentNode{index: 0},
				cmp:   "eq",
				right: literalNode{value: 5},
			},
		},
	})
}

func (s *RulesSuite) Test_parseExpressionWithMultiplication(c *C) {
	result, _ := parseExpression("arg0 == 12 * 3")
	c.Assert(result.String(), Equals, "(eq arg0 (mul 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithAddition(c *C) {
	result, _ := parseExpression("arg0 == 12 + 3")
	c.Assert(result.String(), Equals, "(eq arg0 (add 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithDivision(c *C) {
	result, _ := parseExpression("arg0 == 12 / 3")
	c.Assert(result.String(), Equals, "(eq arg0 (quo 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithSubtraction(c *C) {
	result, _ := parseExpression("arg0 == 12 - 3")
	c.Assert(result.String(), Equals, "(eq arg0 (sub 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryAnd(c *C) {
	result, _ := parseExpression("arg0 == 0 & 1")
	c.Assert(result.String(), Equals, "(eq arg0 (and 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryOr(c *C) {
	result, _ := parseExpression("arg0 == 0 | 1")
	c.Assert(result.String(), Equals, "(eq arg0 (or 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryXor(c *C) {
	result, _ := parseExpression("arg0 == 0 ^ 1")
	c.Assert(result.String(), Equals, "(eq arg0 (xor 0 1))")
}

func (s *RulesSuite) Test_parseAExpressionBinaryNegation(c *C) {
	c.Skip("not yet implemented, check binary negation syntax")
	result, _ := parseExpression("arg0 == ^0")
	c.Assert(result.String(), Equals, "(eq arg0 (bnot 0))")
}

func (s *RulesSuite) Test_parseAExpressionLeftShift(c *C) {
	result, _ := parseExpression("arg0 == 2 << 1")
	c.Assert(result.String(), Equals, "(eq arg0 (shl 2 1))")
}

func (s *RulesSuite) Test_parseAExpressionRightShift(c *C) {
	result, _ := parseExpression("arg0 == 2 >> 1")
	c.Assert(result.String(), Equals, "(eq arg0 (shr 2 1))")
}

func (s *RulesSuite) Test_parseAExpressionWithModulo(c *C) {
	result, _ := parseExpression("arg0 == 12 % 3")
	c.Assert(result.String(), Equals, "(eq arg0 (rem 12 3))")
}

func (s *RulesSuite) Test_parseAExpressionWithBooleanAnd(c *C) {
	// is this the expected syntax for this or should we ever
	// consider cases like arg0 && arg1
	result, _ := parseExpression("arg0 == 0 && arg1 == 0")
	c.Assert(result.String(), Equals, "(land (eq arg0 0) (eq arg1 0))")
}

func (s *RulesSuite) Test_parseAExpressionWithBooleanNegation(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("!arg0")
	c.Assert(result.String(), Equals, "(not arg0)")
}

func (s *RulesSuite) Test_parseAExpressionWithNotEqual(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 != arg1")
	c.Assert(result.String(), Equals, "(neq arg0 arg1")
}

func (s *RulesSuite) Test_parseAExpressionWithGreaterThanOrEqualTo(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 >= arg1")
	c.Assert(result.String(), Equals, "(geq arg0 arg1")
}

func (s *RulesSuite) Test_parseAExpressionWithLessThan(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 < arg1")
	c.Assert(result.String(), Equals, "(lss arg0 arg1")
}

func (s *RulesSuite) Test_parseAExpressionWithLessThanOrEqualTo(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 <= arg1")
	c.Assert(result.String(), Equals, "(leq arg0")
}

func (s *RulesSuite) Test_parseAExpressionWithBitSets(c *C) {
	c.Skip("not yet implemented, check syntax against how we use binary and")
	result, _ := parseExpression("arg0 & val")
	c.Assert(result.String(), Equals, "(set arg0")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusion(c *C) {
	c.Skip("not yet implemented, check syntax about set")
	result, _ := parseExpression("in(arg0, 1, 2)")
	c.Assert(result.String(), Equals, "(in arg0 {1, 2}")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusionLargerSet(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("in(arg0, 1, 2, 3, 4)")
	c.Assert(result.String(), Equals, "(in arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithInclusionWithWhitespace(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("in(arg0, 1,   2,   3,   4)")
	c.Assert(result.String(), Equals, "(in arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithNotInclusion(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("notIn(arg0, 1, 2)")
	c.Assert(result.String(), Equals, "(notIn arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithNotInclusionLargerSet(c *C) {
	c.Skip("not yet implemented, check syntax about set syntax")
	result, _ := parseExpression("notin(arg0, 1, 2, 3, 4)")
	c.Assert(result.String(), Equals, "(notin arg0 {1, 2, 3, 4}")
}

func (s *RulesSuite) Test_parseAExpressionWithTrue(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("true")
	c.Assert(result.String(), Equals, "1")
}

func (s *RulesSuite) Test_parseAExpressionWithFalse(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("false")
	c.Assert(result.String(), Equals, "0")
}

func (s *RulesSuite) Test_parseAExpressionWith0AsFalse(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("0")
	c.Assert(result.String(), Equals, "0")
}

func (s *RulesSuite) Test_parseAExpressionWithNestedOperatorsWithParens(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 == (12 + 3) * 2")
	c.Assert(result.String(), Equals, "(eq arg0 (* (+ 12 3) 2)))")
}

func (s *RulesSuite) Test_parseAExpressionWithNestedOperators(c *C) {
	c.Skip("not yet implemented")
	result, _ := parseExpression("arg0 == 12 + 3 * 2")
	c.Assert(result.String(), Equals, "(eq arg0 (+ 12 (* 3 2)))")
}

func (s *RulesSuite) Test_parseAExpressionWithInvalidArithmeticOperator(c *C) {
	c.Skip("not yet implemented, error handling")
	result, _ := parseExpression("arg0 == 12 _ 3")
	c.Assert(result.String(), Equals, "(eq arg0 (add 12 3))")
}

//	result, _ := doParse("read2: arg0 > 0")

//	c.Assert(result.String(), DeepEquals, "(read2 (gt (arg 0) (literal 0)))")

// func (s *RulesSuite) Test_parsesSlightlyMoreComplicatedRule(c *C) {
// 	result, _ := doParse("write: arg1 == 42 || arg0 + 1 == 15 && (arg3 == 1 || arg4 == 2)")

// 	c.Assert(result, DeepEquals, []rule{
// 		rule{
// 			syscall: "write",
// 			expression: orExpr{
// 				left: equalsComparison{
// 					left: argumentNode{index: 1},
// 					right: literalNode{value: 42},
// 				},
// 				right: andExpr{
// 					left: equalsComparison{
// 						left: addition{
// 							left: argumentNode{index: 0},
// 							right: literalNode{value: 1},
// 						},
// 						right: literalNode{value: 15},
// 					},
// 					right: orExpr{
// 						left: equalsComparison{
// 							left: argumentNode{index: 3},
// 							right: literalNode{value: 1},
// 						},
// 						right: equalsComparison{
// 							left: argumentNode{index: 4},
// 							right: literalNode{value: 2},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// }
