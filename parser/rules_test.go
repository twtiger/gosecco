package main

// import (
// 	"testing"
//   . "gopkg.in/check.v1"
// )

// func Test(t *testing.T) { TestingT(t) }

// type RulesSuite struct{}

// var _ = Suite(&RulesSuite{})

// func (s *RulesSuite) Test_parsesSimpleRule(c *C) {
// 	result, _ := doParse("read: 1")

// 	c.Assert(result, DeepEquals, []rule{
// 		rule{
// 			syscall: "read",
// 			expression: trueLiteral{},
// 		},
// 	})
// }

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
