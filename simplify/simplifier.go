package simplify

import "github.com/twtiger/go-seccomp/tree"

func Simplify(inp tree.Expression) tree.Expression {
	s := &simplifier{inp}
	inp.Accept(s)
	return s.result
}

type simplifier struct {
	result tree.Expression
}

// AcceptAnd implements Visitor
func (*simplifier) AcceptAnd(tree.And) {}

// AcceptArgument implements Visitor
func (*simplifier) AcceptArgument(tree.Argument) {}

// AcceptBinaryNegation implements Visitor
func (*simplifier) AcceptBinaryNegation(tree.BinaryNegation) {}

// AcceptBooleanLiteral implements Visitor
func (*simplifier) AcceptBooleanLiteral(tree.BooleanLiteral) {}

// AcceptCall implements Visitor
func (*simplifier) AcceptCall(tree.Call) {}

// AcceptComparison implements Visitor
func (*simplifier) AcceptComparison(tree.Comparison) {}

// AcceptInclusion implements Visitor
func (*simplifier) AcceptInclusion(tree.Inclusion) {}

// AcceptNegation implements Visitor
func (*simplifier) AcceptNegation(tree.Negation) {}

// AcceptNumericLiteral implements Visitor
func (*simplifier) AcceptNumericLiteral(x tree.NumericLiteral) {}

// AcceptOr implements Visitor
func (*simplifier) AcceptOr(tree.Or) {}

// AcceptVariable implements Visitor
func (*simplifier) AcceptVariable(tree.Variable) {}
