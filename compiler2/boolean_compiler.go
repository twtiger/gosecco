package compiler2

import "github.com/twtiger/gosecco/tree"

// The boolean compiler uses the stack for simplicity, but we could probably do without
// It generates suboptimal code, expecting a peephole stage after
// It will always take jump points as arguments
// Jump points are arbitary types that represents where to jump.
// All the different boolean situations can be represented using this structure.
// The conditional compiler is also a boolean compiler

type booleanCompilerVisitor struct {
	ctx      *compilerContext
	err      error
	topLevel bool
	jt       label
	jf       label
}

func compileBoolean(ctx *compilerContext, inp tree.Expression, topLevel bool, jt, jf label) error {
	v := &booleanCompilerVisitor{ctx: ctx, jt: jt, jf: jf, topLevel: topLevel}
	inp.Accept(v)
	return v.err
}

func (s *booleanCompilerVisitor) AcceptAnd(v tree.And) {
	// TODO: errors here
	next := s.ctx.newLabel()
	compileBoolean(s.ctx, v.Left, false, next, s.jf)
	s.ctx.labelHere(next)
	compileBoolean(s.ctx, v.Right, false, s.jt, s.jf)
}

// AcceptArgument implements Visitor
func (s *booleanCompilerVisitor) AcceptArgument(v tree.Argument) {
	panic("XXX: generate error here")
}

// AcceptArithmetic implements Visitor
func (s *booleanCompilerVisitor) AcceptArithmetic(v tree.Arithmetic) {
	panic("XXX: generate error here")
}

// AcceptBinaryNegation implements Visitor
func (s *booleanCompilerVisitor) AcceptBinaryNegation(v tree.BinaryNegation) {
	panic("XXX: generate error here")
}

// AcceptBooleanLiteral implements Visitor
func (s *booleanCompilerVisitor) AcceptBooleanLiteral(v tree.BooleanLiteral) {
	if s.topLevel {
		// This should probably just jump to success
		// TODO: compile here
	} else {
		panic("XXX: generate error here")
	}
}

// AcceptCall implements Visitor
func (s *booleanCompilerVisitor) AcceptCall(v tree.Call) {
	panic("XXX: generate error here")
}

// AcceptInclusion implements Visitor
func (s *booleanCompilerVisitor) AcceptInclusion(v tree.Inclusion) {
	panic("XXX: generate error here")
}

// AcceptNegation implements Visitor
func (s *booleanCompilerVisitor) AcceptNegation(v tree.Negation) {
	s.err = compileBoolean(s.ctx, v.Operand, false, s.jf, s.jt)
}

// AcceptNumericLiteral implements Visitor
func (s *booleanCompilerVisitor) AcceptNumericLiteral(v tree.NumericLiteral) {
	panic("XXX: generate error here")
}

// AcceptOr implements Visitor
func (s *booleanCompilerVisitor) AcceptOr(v tree.Or) {
	// TODO: errors
	next := s.ctx.newLabel()
	compileBoolean(s.ctx, v.Left, false, s.jt, next)
	s.ctx.labelHere(next)
	compileBoolean(s.ctx, v.Right, false, s.jt, s.jf)
}

// AcceptVariable implements Visitor
func (s *booleanCompilerVisitor) AcceptVariable(v tree.Variable) {
	panic("XXX: generate error here")
}
