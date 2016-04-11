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

	c.Assert(result.expression.String(), Equals, "(or (eq arg0 4) (eq arg0 5))")
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

func (s *RulesSuite) Test_parseARuleWithArithmetic(c *C) {
	result, _ := parseRule("read5: arg0 == 12 * 3")
	c.Assert(result.expression.String(), Equals, "(eq arg0 (* 12 3))")
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
