package unifier

import "github.com/twtiger/go-seccomp/tree"

// Unify variables within rules from macros
func Unify(r tree.RawPolicy) tree.Policy {
	var rules []tree.Rule
	macros := make(map[string]tree.Macro)
	for _, e := range r.RuleOrMacros {
		switch v := e.(type) {
		case tree.Rule:
			rules = append(rules, replaceFreeNames(v, macros))
		case tree.Macro:
			macros[v.Name] = v
		}
	}
	return tree.Policy{Macros: macros, Rules: rules}
}

func replaceFreeNames(r tree.Rule, macros map[string]tree.Macro) tree.Rule {
	return tree.Rule{
		Name:           r.Name,
		PositiveAction: r.PositiveAction,
		NegativeAction: r.NegativeAction,
		Body:           replace(r.Body, macros),
	}
}

func replace(x tree.Expression, macros map[string]tree.Macro) tree.Expression {
	r := &replacer{x, macros}
	x.Accept(r)
	return r.expression
}
