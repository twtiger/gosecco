package parser

import . "gopkg.in/check.v1"

type BasicsSuite struct{}

var _ = Suite(&BasicsSuite{})

func (s *BasicsSuite) Test_trueLiteral_String(c *C) {
	c.Assert(trueLiteral{}.String(), Equals, "1")
}

func (s *BasicsSuite) Test_literalNode_String(c *C) {
	c.Assert(literalNode{42}.String(), Equals, "42")
}

func (s *BasicsSuite) Test_argumentNode_String(c *C) {
	c.Assert(argumentNode{2}.String(), Equals, "arg2")
}

func (s *BasicsSuite) Test_trueLiteral_Repr(c *C) {
	c.Assert(trueLiteral{}.Repr(), Equals, "1")
}

func (s *BasicsSuite) Test_literalNode_Repr(c *C) {
	c.Assert(literalNode{42}.Repr(), Equals, "42")
}

func (s *BasicsSuite) Test_argumentNode_Repr(c *C) {
	c.Assert(argumentNode{2}.Repr(), Equals, "arg2")
}

func (s *BasicsSuite) Test_arithmetic_String(c *C) {
	left := literalNode{42}
	right := literalNode{15}
	op := "+"
	c.Assert(arithmetic{left, right, op}.String(), Equals, "(+ 42 15)")
}

func (s *BasicsSuite) Test_arithmetic_Repr(c *C) {
	left := argumentNode{2}
	right := literalNode{15}
	op := "+"
	c.Assert(arithmetic{left, right, op}.Repr(), Equals, "arg2 + 15")
}

func (s *BasicsSuite) Test_equalsComparison_String(c *C) {
	c.Assert(equalsComparison{literalNode{42}, literalNode{15}}.String(), Equals, "(== 42 15)")
}

func (s *BasicsSuite) Test_equalsComparison_Repr(c *C) {
	c.Assert(equalsComparison{argumentNode{2}, literalNode{15}}.Repr(), Equals, "arg2 == 15")
}

func (s *BasicsSuite) Test_orExpr_String(c *C) {
	c.Assert(orExpr{equalsComparison{argumentNode{2}, literalNode{15}}, trueLiteral{}}.String(), Equals, "(lor (== arg2 15) 1)")
}

func (s *BasicsSuite) Test_orExpr_Repr(c *C) {
	c.Assert(orExpr{equalsComparison{argumentNode{2}, literalNode{15}}, trueLiteral{}}.Repr(), Equals, "arg2 == 15 || 1")
}

func (s *BasicsSuite) Test_andExpr_String(c *C) {
	c.Assert(andExpr{equalsComparison{argumentNode{2}, literalNode{15}}, trueLiteral{}}.String(), Equals, "(land (== arg2 15) 1)")
}

func (s *BasicsSuite) Test_andExpr_Repr(c *C) {
	c.Assert(andExpr{equalsComparison{argumentNode{2}, literalNode{15}}, trueLiteral{}}.Repr(), Equals, "arg2 == 15 && 1")
}
