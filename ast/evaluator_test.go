package ast

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type EvaluatorSuite struct{}

var _ = Suite(&EvaluatorSuite{})

func (s *EvaluatorSuite) Test_bla(c *C) {
	eval := &EvaluatorVisitor{}
	a := Arithmetic{Op: MINUS, Left: NumericLiteral{42}, Right: NumericLiteral{23}}
	a.Accept(eval)

	fmt.Printf("Result: %d\n", eval.popNumeric())
}
