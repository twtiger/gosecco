package parser

import (
	"regexp"
	"strings"
)

type ruleHead struct {
	syscall  string
	positive string
	negative string
}

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

func parseRuleHead(s string) (ruleHead, bool) {
	match := ruleHeadRE.FindStringSubmatch(s)
	if match != nil {
		positive, negative, ok := findPositiveAndNegative(strings.Split(match[2], ","))
		return ruleHead{syscall: match[1], positive: positive, negative: negative}, ok
	}
	return ruleHead{}, false
}
