package tree

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type EvaluatorSuite struct{}

var _ = Suite(&EvaluatorSuite{})

func (s *EvaluatorSuite) Test_simpleArithmetic(c *C) {
	eval := &EvaluatorVisitor{}
	a := Arithmetic{Op: MINUS, Left: NumericLiteral{Value: 42}, Right: NumericLiteral{Value: 23}}
	a.Accept(eval)
	c.Assert(eval.popNumeric(), Equals, uint32(0x13))
}

func (s *EvaluatorSuite) Test_complicatedArithmetic(c *C) {
	eval := &EvaluatorVisitor{}
	a := Arithmetic{
		Left:  Arithmetic{Op: PLUS, Left: NumericLiteral{Value: 42}, Right: BinaryNegation{NumericLiteral{Value: 23}}},
		Op:    MULT,
		Right: Arithmetic{Op: DIV, Left: NumericLiteral{Value: 4}, Right: NumericLiteral{Value: 2}},
	}

	a.Accept(eval)
	c.Assert(eval.popNumeric(), Equals, uint32(0x24))
}

func (s *EvaluatorSuite) Test_boolean(c *C) {
	eval := &EvaluatorVisitor{}
	a := Negation{And{
		Left: BooleanLiteral{true},
		Right: Or{
			Left:  Comparison{Op: EQL, Left: NumericLiteral{Value: 42}, Right: NumericLiteral{Value: 23}},
			Right: Comparison{Op: NEQL, Left: NumericLiteral{Value: 42}, Right: NumericLiteral{Value: 42}},
		},
	}}
	a.Accept(eval)
	c.Assert(eval.popBoolean(), Equals, true)
}
