package tree

// Visitor is a visitor for all parse nodes
type Visitor interface {
	AcceptAnd(And)
	AcceptArgument(Argument)
	AcceptArithmetic(Arithmetic)
	AcceptBinaryNegation(BinaryNegation)
	AcceptBooleanLiteral(BooleanLiteral)
	AcceptCall(Call)
	AcceptComparison(Comparison)
	AcceptInclusion(Inclusion)
	AcceptNegation(Negation)
	AcceptNumericLiteral(NumericLiteral)
	AcceptOr(Or)
	AcceptVariable(Variable)
}

// Expression is an AST expression
type Expression interface {
	Accept(Visitor)
}

// EmptyVisitor provides empty default methods for visitors that only care about one type
type EmptyVisitor struct{}

// AcceptAnd implements Visitor
func (*EmptyVisitor) AcceptAnd(And) {}

// AcceptArgument implements Visitor
func (*EmptyVisitor) AcceptArgument(Argument) {}

// AcceptArithmetic implements Visitor
func (*EmptyVisitor) AcceptArithmetic(Arithmetic) {}

// AcceptBinaryNegation implements Visitor
func (*EmptyVisitor) AcceptBinaryNegation(BinaryNegation) {}

// AcceptBooleanLiteral implements Visitor
func (*EmptyVisitor) AcceptBooleanLiteral(BooleanLiteral) {}

// AcceptCall implements Visitor
func (*EmptyVisitor) AcceptCall(Call) {}

// AcceptComparison implements Visitor
func (*EmptyVisitor) AcceptComparison(Comparison) {}

// AcceptInclusion implements Visitor
func (*EmptyVisitor) AcceptInclusion(Inclusion) {}

// AcceptNegation implements Visitor
func (*EmptyVisitor) AcceptNegation(Negation) {}

// AcceptNumericLiteral implements Visitor
func (*EmptyVisitor) AcceptNumericLiteral(NumericLiteral) {}

// AcceptOr implements Visitor
func (*EmptyVisitor) AcceptOr(Or) {}

// AcceptVariable implements Visitor
func (*EmptyVisitor) AcceptVariable(Variable) {}
