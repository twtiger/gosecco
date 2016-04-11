package parser

import "fmt"

type expression interface {
	String() string
	Repr() string
}

type trueLiteral struct{}

func (trueLiteral) String() string {
	return "1"
}

func (trueLiteral) Repr() string {
	return "1"
}

type rule struct {
	syscall    string
	expression expression
}

type booleanExpression interface {
	String() string
	Repr() string
}

type integerExpression interface {
	String() string
	Repr() string
}

type orExpr struct {
	left, right booleanExpression
}

func (a orExpr) String() string {
	return fmt.Sprintf("(or %s %s)", a.left.String(), a.right.String())
}

func (a orExpr) Repr() string {
	return fmt.Sprintf("%s || %s", a.left.Repr(), a.right.Repr())
}

type andExpr struct {
	left, right booleanExpression
}

func (a andExpr) String() string {
	return fmt.Sprintf("(and %s %s)", a.left.String(), a.right.String())
}

func (a andExpr) Repr() string {
	return fmt.Sprintf("%s && %s", a.left.Repr(), a.right.Repr())
}

type equalsComparison struct {
	left, right integerExpression
}

func (e equalsComparison) String() string {
	return fmt.Sprintf("(== %s %s)", e.left.String(), e.right.String())
}

func (e equalsComparison) Repr() string {
	return fmt.Sprintf("%s == %s", e.left.Repr(), e.right.Repr())
}

type argumentNode struct {
	index int
}

func (a argumentNode) String() string {
	return fmt.Sprintf("arg%d", a.index)
}

func (a argumentNode) Repr() string {
	return a.String()
}

type literalNode struct {
	value int
}

func (l literalNode) String() string {
	return fmt.Sprintf("%d", l.value)
}

func (l literalNode) Repr() string {
	return fmt.Sprintf("%d", l.value)
}

type comparison struct {
	left, right integerExpression
	cmp         string
}

func (c comparison) String() string {
	return fmt.Sprintf("(%s %s %s)", c.cmp, c.left.String(), c.right.String())
}

func (c comparison) Repr() string {
	return fmt.Sprintf("%s %s %s", c.left.Repr(), c.cmp, c.right.Repr())
}

type arithmetic struct {
	left, right integerExpression
	op          string
}

func (a arithmetic) String() string {
	return fmt.Sprintf("(%s %s %s)", a.op, a.left.String(), a.right.String())
}

func (a arithmetic) Repr() string {
	return fmt.Sprintf("%s %s %s", a.left.Repr(), a.op, a.right.Repr())
}
