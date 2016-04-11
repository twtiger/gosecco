package main

type expression interface{}

type trueLiteral struct{}

type rule struct {
	syscall    string
	expression expression
}

func doParse(source string) ([]rule, error) {
	return nil, nil
}

type BooleanExpression interface {
	EvaluateBool() bool
}

type IntegerExpression interface {
	EvaluateInt() int
}

type orExpr struct {
	left, right BooleanExpression
}
type andExpr struct {
	left, right BooleanExpression
}
type equalsComparison struct {
	left, right IntegerExpression
}
type argumentNode struct {
	index int
}
type literalNode struct {
	value int
}
type addition struct {
	left, right IntegerExpression
}

func (e orExpr) EvaluateBool() bool {
	return e.left.EvaluateBool() || e.right.EvaluateBool()
}

func (e andExpr) EvaluateBool() bool {
	if e.left.EvaluateBool() {
		return e.right.EvaluateBool()
	}
	return false
}

func (e equalsComparison) EvaluateBool() bool {
	return e.left.EvaluateInt() == e.right.EvaluateInt()
}

func (e addition) EvaluateInt() int {
	return e.left.EvaluateInt() + e.right.EvaluateInt()
}

func (e argumentNode) EvaluateInt() int {
	return 0
}

func (e literalNode) EvaluateInt() int {
	return e.value
}
