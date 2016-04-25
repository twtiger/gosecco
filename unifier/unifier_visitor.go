package unifier

import (
	"errors"

	"github.com/twtiger/gosecco/tree"
)

type replacer struct {
	expression tree.Expression
	macros     map[string]tree.Macro
	err        error
}

func (r *replacer) AcceptAnd(b tree.And) {
	var left tree.Boolean
	var right tree.Boolean
	left, r.err = replace(b.Left, r.macros)
	if r.err == nil {
		right, r.err = replace(b.Right, r.macros)
		r.expression = tree.And{Left: left, Right: right}
	}
}

func (r *replacer) AcceptArgument(tree.Argument) {}

func (r *replacer) AcceptArithmetic(b tree.Arithmetic) {
	var left tree.Numeric
	var right tree.Numeric
	left, r.err = replace(b.Left, r.macros)
	if r.err == nil {
		right, r.err = replace(b.Right, r.macros)
		r.expression = tree.Arithmetic{Left: left, Op: b.Op, Right: right}
	}
}

func (r *replacer) AcceptBinaryNegation(b tree.BinaryNegation) {
	var op tree.Numeric
	op, r.err = replace(b.Operand, r.macros)
	r.expression = tree.BinaryNegation{op}
}

func (r *replacer) AcceptBooleanLiteral(tree.BooleanLiteral) {}

func (r *replacer) AcceptCall(b tree.Call) {
	v := r.macros[b.Name] // we get the name of the macro

	// write: call(arg0)

	// call becomes our macro name
	// then we make macros like: arg0: x
	// then we add it to our list of macros
	// then we reduce it

	nm := make(map[string]tree.Macro)
	for i, e := range b.Args {
		m := tree.Macro{Name: v.ArgumentNames[i], Body: e} // we make a new macro
		nm[v.ArgumentNames[i]] = m
	}

	exp, err := replace(v.Body, nm)
	r.err = err

	for k, v := range r.macros {
		nm[k] = v
	}

	if r.err == nil {
		r.expression, r.err = replace(exp, nm)
	}
}

func (r *replacer) AcceptComparison(b tree.Comparison) {
	var left tree.Numeric
	var right tree.Numeric

	left, r.err = replace(b.Left, r.macros)

	if r.err == nil {
		right, r.err = replace(b.Right, r.macros)
		r.expression = tree.Comparison{
			Left:  left,
			Op:    b.Op,
			Right: right,
		}
	}

}

func (r *replacer) AcceptInclusion(b tree.Inclusion) {
	var rights []tree.Numeric
	for _, e := range b.Rights {
		right, err := replace(e, r.macros)
		if err != nil {
			r.err = err
		}
		rights = append(rights, right)
	}
	left, err := replace(b.Left, r.macros)
	if err != nil {
		r.err = err
	}

	r.expression = tree.Inclusion{Positive: b.Positive, Left: left, Rights: rights}
}

func (r *replacer) AcceptNegation(b tree.Negation) {
	var op tree.Numeric
	op, r.err = replace(b.Operand, r.macros)

	r.expression = tree.Negation{Operand: op}
}

func (r *replacer) AcceptNumericLiteral(tree.NumericLiteral) {}

func (r *replacer) AcceptOr(b tree.Or) {
	var left tree.Boolean
	var right tree.Boolean

	left, r.err = replace(b.Left, r.macros)

	if r.err == nil {
		right, r.err = replace(b.Right, r.macros)
		r.expression = tree.And{Left: left, Right: right}
		r.expression = tree.Or{Left: left, Right: right}
	}
}

func (r *replacer) AcceptVariable(b tree.Variable) {
	// TODO handle case where a macro is called that should take a variable but is not given one
	expr, ok := r.macros[b.Name]
	if ok {
		r.expression = expr.Body
	} else {
		r.err = errors.New("Variable not defined")
	}
}
