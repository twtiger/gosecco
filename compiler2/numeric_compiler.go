package compiler2

import "github.com/twtiger/gosecco/tree"

// arithmeticCompilerVisitor is instantiated and run once we hit a numeric expression
// the end result will be the full byte code for the numeric expression.
// it will generate an error for anything else
type numericCompilerVisitor struct {
	ctx *compilerContext
}

func compileNumeric(ctx *compilerContext, inp tree.Expression) {
	inp.Accept(&numericCompilerVisitor{ctx})
}

func (s *numericCompilerVisitor) AcceptAnd(v tree.And) {
	panic("TODO: generate error here")
}

// AcceptArgument implements Visitor
func (s *numericCompilerVisitor) AcceptArgument(v tree.Argument) {
	// TODO: implement
}

// AcceptArithmetic implements Visitor
func (s *numericCompilerVisitor) AcceptArithmetic(v tree.Arithmetic) {
	// TODO: implement
}

// AcceptBinaryNegation implements Visitor
func (s *numericCompilerVisitor) AcceptBinaryNegation(v tree.BinaryNegation) {
	// TODO: implement
}

// AcceptBooleanLiteral implements Visitor
func (s *numericCompilerVisitor) AcceptBooleanLiteral(v tree.BooleanLiteral) {
	panic("TODO: generate error here")
}

// AcceptCall implements Visitor
func (s *numericCompilerVisitor) AcceptCall(v tree.Call) {
	panic("TODO: generate error here")
}

// AcceptComparison implements Visitor
func (s *numericCompilerVisitor) AcceptComparison(v tree.Comparison) {
	panic("TODO: generate error here")
}

// AcceptInclusion implements Visitor
func (s *numericCompilerVisitor) AcceptInclusion(v tree.Inclusion) {
	panic("TODO: generate error here")
}

// AcceptNegation implements Visitor
func (s *numericCompilerVisitor) AcceptNegation(v tree.Negation) {
	panic("TODO: generate error here")
}

// AcceptNumericLiteral implements Visitor
func (s *numericCompilerVisitor) AcceptNumericLiteral(v tree.NumericLiteral) {
	s.ctx.op(OP_LOAD_VAL, uint32(v.Value))
}

// AcceptOr implements Visitor
func (s *numericCompilerVisitor) AcceptOr(v tree.Or) {
	panic("TODO: generate error here")
}

// AcceptVariable implements Visitor
func (s *numericCompilerVisitor) AcceptVariable(v tree.Variable) {
	panic("TODO: generate error here")
}
