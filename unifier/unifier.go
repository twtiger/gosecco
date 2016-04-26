package unifier

import (
	"strconv"

	"github.com/twtiger/gosecco/tree"
)

func getDefaultActions(t tree.Macro) string {
	var def string
	switch f := t.Body.(type) {
	case tree.NumericLiteral:
		def = strconv.Itoa(int(f.Value))
	case tree.Variable:
		def = f.Name
	}
	return def
}

func setDefaultActions(defaults []tree.Macro, pt tree.PolicyType, enforce bool) (string, string) {
	var defpos string
	var defneg string

	for _, v := range defaults {
		if v.Name == "DEFAULT_POSITIVE" {
			defpos = getDefaultActions(v)
		} else if v.Name == "DEFAULT_NEGATIVE" {
			defneg = getDefaultActions(v)
		}
	}

	if len(defpos) == 0 {
		if enforce == true {
			if pt == tree.WhiteList {
				defpos = "allow"
			} else { // only other option is Blacklist
				defpos = "kill"
			}
		} else {
			if pt == tree.WhiteList {
				defpos = "allow"
			} else {
				defpos = "trace"
			}
		}
	}

	if len(defneg) == 0 {
		if enforce == true {
			if pt == tree.WhiteList {
				defneg = "kill"
			} else {
				defneg = "allow"
			}
		} else {
			if pt == tree.WhiteList {
				defneg = "trace"
			} else {
				defneg = "allow"
			}
		}
	}

	return defpos, defneg
}

// Unify variables within rules from macros
func Unify(r tree.RawPolicy, enforce bool) (tree.Policy, error) {
	var err error
	var rules []tree.Rule
	macros := make(map[string]tree.Macro)
	defs := make([]tree.Macro, 0)
	for _, e := range r.RuleOrMacros {
		switch v := e.(type) {
		case tree.Rule:
			var r tree.Rule
			r, err = replaceFreeNames(v, macros)
			rules = append(rules, r)
		case tree.Macro:
			if v.Name != "DEFAULT_POSITIVE" && v.Name != "DEFAULT_NEGATIVE" {
				macros[v.Name] = v
			} else {
				defs = append(defs, v)
			}
		}
	}
	defpos, defneg := setDefaultActions(defs, r.ListType, enforce)
	return tree.Policy{DefaultPositiveAction: defpos, DefaultNegativeAction: defneg, Macros: macros, Rules: rules}, err
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
