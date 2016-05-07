package simplifier

import "github.com/twtiger/gosecco/tree"

// AcceptArgument implements Visitor
func (*simplifier) AcceptArgument(tree.Argument) {}

// AcceptBooleanLiteral implements Visitor
func (*simplifier) AcceptBooleanLiteral(tree.BooleanLiteral) {}

// AcceptNumericLiteral implements Visitor
func (*simplifier) AcceptNumericLiteral(x tree.NumericLiteral) {}

// AcceptVariable implements Visitor
func (*simplifier) AcceptVariable(tree.Variable) {}
