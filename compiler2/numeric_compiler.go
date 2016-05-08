package compiler2

import (
	"errors"

	"github.com/twtiger/gosecco/tree"
)

// The numeric compiler actually needs a stack to function
// The ctx will contain a number - this is the current stack top
// We will use the  BPF_LD+BPF_MEM and BPF_ST instructions for this
// The value in context will contain the next free stack entry.
// This stack can contain max BPF_MEMWORDS entries, and will generate an error if we try to use more than that
// In general we will first evaluate Right, then push that on the stack, evaluate Left, then execute the action.

// arithmeticCompilerVisitor is instantiated and run once we hit a numeric expression
// the end result will be the full byte code for the numeric expression.
// it will generate an error for anything else
type numericCompilerVisitor struct {
	ctx *compilerContext
	err error
}

func compileNumeric(ctx *compilerContext, inp tree.Expression) error {
	v := &numericCompilerVisitor{ctx: ctx}
	inp.Accept(v)
	return v.err
}

func (s *numericCompilerVisitor) AcceptAnd(v tree.And) {
	panic("XXX: generate error here")
}

const argumentsStartIndex = uint32(0x10)

// AcceptArgument implements Visitor
func (s *numericCompilerVisitor) AcceptArgument(v tree.Argument) {
	argIndex := argumentsStartIndex + uint32(v.Index*8)
	if v.Type == tree.Low {
		argIndex += 4
	}

	s.ctx.op(OP_LOAD, argIndex)
}

var arithOps = map[tree.ArithmeticType]uint16{
	tree.PLUS:   OP_ADD_X,
	tree.MINUS:  OP_SUB_X,
	tree.MULT:   OP_MUL_X,
	tree.DIV:    OP_DIV_X,
	tree.BINAND: OP_AND_X,
	tree.BINOR:  OP_OR_X,
	tree.BINXOR: OP_XOR_X,
	tree.LSH:    OP_LSH_X,
	tree.RSH:    OP_RSH_X,
	tree.MOD:    OP_MOD_X,
}

// AcceptArithmetic implements Visitor
func (s *numericCompilerVisitor) AcceptArithmetic(v tree.Arithmetic) {
	if err := compileNumeric(s.ctx, v.Right); err != nil {
		s.err = err
		return
	}

	if err := s.ctx.pushAToStack(); err != nil {
		s.err = err
		return
	}

	if err := compileNumeric(s.ctx, v.Left); err != nil {
		s.err = err
		return
	}

	if err := s.ctx.popStackToX(); err != nil {
		s.err = err
		return
	}

	arithOp, ok := arithOps[v.Op]
	if !ok {
		s.err = errors.New("BLA")
		return
	}
	s.ctx.op(arithOp, 0)
}

// AcceptBinaryNegation implements Visitor
func (s *numericCompilerVisitor) AcceptBinaryNegation(v tree.BinaryNegation) {
	panic("XXX: generate error here")
}

// AcceptBooleanLiteral implements Visitor
func (s *numericCompilerVisitor) AcceptBooleanLiteral(v tree.BooleanLiteral) {
	panic("XXX: generate error here")
}

// AcceptCall implements Visitor
func (s *numericCompilerVisitor) AcceptCall(v tree.Call) {
	panic("XXX: generate error here")
}

// AcceptComparison implements Visitor
func (s *numericCompilerVisitor) AcceptComparison(v tree.Comparison) {
	panic("XXX: generate error here")
}

// AcceptInclusion implements Visitor
func (s *numericCompilerVisitor) AcceptInclusion(v tree.Inclusion) {
	panic("XXX: generate error here")
}

// AcceptNegation implements Visitor
func (s *numericCompilerVisitor) AcceptNegation(v tree.Negation) {
	panic("XXX: generate error here")
}

// AcceptNumericLiteral implements Visitor
func (s *numericCompilerVisitor) AcceptNumericLiteral(v tree.NumericLiteral) {
	s.ctx.op(OP_LOAD_VAL, uint32(v.Value))
}

// AcceptOr implements Visitor
func (s *numericCompilerVisitor) AcceptOr(v tree.Or) {
	panic("XXX: generate error here")
}

// AcceptVariable implements Visitor
func (s *numericCompilerVisitor) AcceptVariable(v tree.Variable) {
	panic("XXX: generate error here")
}
