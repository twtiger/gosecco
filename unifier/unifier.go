package unifier

import "github.com/twtiger/go-seccomp/tree"

// Unify variables within rules from macros
func Unify(r tree.RawPolicy) (tree.Policy, error) {
	var err error
	var rules []tree.Rule
	macros := make(map[string]tree.Macro)
	for _, e := range r.RuleOrMacros {
		switch v := e.(type) {
		case tree.Rule:
			var r tree.Rule
			r, err = replaceFreeNames(v, macros)
			rules = append(rules, r)
		case tree.Macro:
			macros[v.Name] = v
		}
	}
	return tree.Policy{Macros: macros, Rules: rules}, err
}

func replaceFreeNames(r tree.Rule, macros map[string]tree.Macro) (tree.Rule, error) {
	body, err := replace(r.Body, macros)
	rule := tree.Rule{
		Name:           r.Name,
		PositiveAction: r.PositiveAction,
		NegativeAction: r.NegativeAction,
		Body:           body,
	}
	return rule, err
}

func replace(x tree.Expression, macros map[string]tree.Macro) (tree.Expression, error) {
	r := &replacer{expression: x, macros: macros, err: nil}
	x.Accept(r)
	if r.err != nil {
		return nil, r.err
	}
	return r.expression, nil
}
