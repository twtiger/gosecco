package unifier

import (
	"fmt"
	"github.com/twtiger/go-seccomp/tree"
)

type replacer struct {
	expression tree.Expression
	macros     map[string]tree.Macro
}

func (r *replacer) AcceptAnd(b tree.And) {
	r.expression = tree.And{Left: replace(b.Left, r.macros), Right: replace(b.Right, r.macros)}
}

func (r *replacer) AcceptArgument(tree.Argument) {}

func (r *replacer) AcceptArithmetic(b tree.Arithmetic) {
	r.expression = tree.Arithmetic{Left: replace(b.Left, r.macros), Op: b.Op, Right: replace(b.Right, r.macros)}
}

func (r *replacer) AcceptBinaryNegation(b tree.BinaryNegation) {
	r.expression = tree.BinaryNegation{replace(b.Operand, r.macros)}
}

func (r *replacer) AcceptBooleanLiteral(tree.BooleanLiteral) {}

func (r *replacer) AcceptCall(b tree.Call) {
	// should be generic
	v := r.macros[b.Name]

	nm := make(map[string]tree.Macro)
	for i, e := range b.Args {
		m := tree.Macro{Name: v.ArgumentNames[i], Body: e}
		nm[v.ArgumentNames[i]] = m
	}

	exp := replace(v.Body, nm)

	for k, v := range r.macros {
		nm[k] = v
	}

	r.expression = replace(exp, nm)
}

func (r *replacer) AcceptComparison(b tree.Comparison) {
	r.expression = tree.Comparison{
		Left:  replace(b.Left, r.macros),
		Op:    b.Op,
		Right: replace(b.Right, r.macros),
	}
}

func (r *replacer) AcceptInclusion(b tree.Inclusion) {
	var rights []tree.Numeric
	for _, e := range b.Rights {
		rights = append(rights, replace(e, r.macros))
	}
	r.expression = tree.Inclusion{Positive: b.Positive, Left: replace(b.Left, r.macros), Rights: rights}
}

func (r *replacer) AcceptNegation(b tree.Negation) {
	r.expression = tree.Negation{Operand: replace(b.Operand, r.macros)}
}

func (r *replacer) AcceptNumericLiteral(tree.NumericLiteral) {}

func (r *replacer) AcceptOr(b tree.Or) {
	r.expression = tree.Or{Left: replace(b.Left, r.macros), Right: replace(b.Right, r.macros)}
}

func (r *replacer) AcceptVariable(b tree.Variable) {
	// TODO: handle missing macro definition
	// TODO: handle case where variable references macro that takes arguments
	expr, ok := r.macros[b.Name]
	if ok {
		r.expression = expr.Body
	} else {
		panic(fmt.Sprintf("sadness: %#v")) //TODO handle this case
	}
}
