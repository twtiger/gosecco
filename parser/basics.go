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

func doParse(source string) ([]rule, error) {
	return nil, nil
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
	return l.String()
}

type addition struct {
	left, right integerExpression
}

func (a addition) String() string {
	return fmt.Sprintf("(+ %s %s)", a.left.String(), a.right.String())
}

func (a addition) Repr() string {
	return fmt.Sprintf("%s + %s", a.left.Repr(), a.right.Repr())
}
