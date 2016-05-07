package simplifier

import "github.com/twtiger/gosecco/tree"

// nullSimplifier does nothing - it returns the same tree as given
// It can be useful as the base for other simplifiers
type nullSimplifier struct {
	result   tree.Expression
	realSelf Simplifier
}

// Simplify implements simplifier
func (s *nullSimplifier) Simplify(inp tree.Expression) tree.Expression {
	inp.Accept(s.realSelf)
	return s.result
}

// AcceptAnd implements Visitor
func (s *nullSimplifier) AcceptAnd(v tree.And) {
	s.result = tree.And{
		Left:  s.realSelf.Simplify(v.Left),
		Right: s.realSelf.Simplify(v.Right),
	}
}

// AcceptArgument implements Visitor
func (s *nullSimplifier) AcceptArgument(v tree.Argument) {
	s.result = v
}

// AcceptArithmetic implements Visitor
func (s *nullSimplifier) AcceptArithmetic(v tree.Arithmetic) {
	s.result = tree.Arithmetic{
		Op:    v.Op,
		Left:  s.realSelf.Simplify(v.Left),
		Right: s.realSelf.Simplify(v.Right),
	}
}

// AcceptBinaryNegation implements Visitor
func (s *nullSimplifier) AcceptBinaryNegation(v tree.BinaryNegation) {
	s.result = tree.BinaryNegation{s.realSelf.Simplify(v.Operand)}
}

// AcceptBooleanLiteral implements Visitor
func (s *nullSimplifier) AcceptBooleanLiteral(v tree.BooleanLiteral) {
	s.result = v
}

// AcceptCall implements Visitor
func (s *nullSimplifier) AcceptCall(v tree.Call) {
	result := make([]tree.Any, len(v.Args))
	for ix, v2 := range v.Args {
		result[ix] = s.realSelf.Simplify(v2)
	}
	s.result = tree.Call{Name: v.Name, Args: result}
}

// AcceptComparison implements Visitor
func (s *nullSimplifier) AcceptComparison(v tree.Comparison) {
	s.result = tree.Comparison{
		Op:    v.Op,
		Left:  s.realSelf.Simplify(v.Left),
		Right: s.realSelf.Simplify(v.Right),
	}
}

// AcceptInclusion implements Visitor
func (s *nullSimplifier) AcceptInclusion(v tree.Inclusion) {
	result := make([]tree.Numeric, len(v.Rights))
	for ix, v2 := range v.Rights {
		result[ix] = s.realSelf.Simplify(v2)
	}
	s.result = tree.Inclusion{
		Positive: v.Positive,
		Left:     s.realSelf.Simplify(v.Left),
		Rights:   result}
}

// AcceptNegation implements Visitor
func (s *nullSimplifier) AcceptNegation(v tree.Negation) {
	s.result = tree.Negation{s.realSelf.Simplify(v.Operand)}
}

// AcceptNumericLiteral implements Visitor
func (s *nullSimplifier) AcceptNumericLiteral(v tree.NumericLiteral) {
	s.result = v
}

// AcceptOr implements Visitor
func (s *nullSimplifier) AcceptOr(v tree.Or) {
	s.result = tree.Or{
		Left:  s.realSelf.Simplify(v.Left),
		Right: s.realSelf.Simplify(v.Right),
	}
}

// AcceptVariable implements Visitor
func (s *nullSimplifier) AcceptVariable(v tree.Variable) {
	s.result = v
}
