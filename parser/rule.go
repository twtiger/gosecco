package parser

import (
	"regexp"
	"strings"

	"github.com/twtiger/go-seccomp/tree"
)

var ruleHeadRE = regexp.MustCompile(`^[[:space:]]*([[:word:]]+)[[:space:]]*(?:\[(.*)\])?[[:space:]]*$`)

func findPositiveAndNegative(ss []string) (string, string, bool) {
	neg, pos := "", ""
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			if strings.HasPrefix(s, "+") {
				if pos != "" {
					return "", "", false
				}
				pos = strings.TrimPrefix(s, "+")
			} else if strings.HasPrefix(s, "-") {
				if neg != "" {
					return "", "", false
				}
				neg = strings.TrimPrefix(s, "-")
			} else {
				return "", "", false
			}
		}
	}
	return pos, neg, true
}

func parseRuleHead(s string) (tree.Rule, bool) {
	match := ruleHeadRE.FindStringSubmatch(s)
	if match != nil {
		positive, negative, ok := findPositiveAndNegative(strings.Split(match[2], ","))
		return tree.Rule{Name: match[1], PositiveAction: positive, NegativeAction: negative}, ok
	}
	return tree.Rule{}, false
}
