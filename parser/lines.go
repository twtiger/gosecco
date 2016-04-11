package parser

import "strings"

type LineType int

const (
	unknownLine LineType = iota
	ruleLine
	commentLine
	assignmentLine
	defaultAssignmentLine
)

func isComment(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "#")
}

func isRule(s string) bool {
	return len(strings.SplitN(s, ":", 2)) == 2
}

func isDefaultAssignment(s string) bool {
	result := strings.SplitN(s, "=", 2)
	if len(result) == 2 {
		c := strings.TrimSpace(result[0])
		return c == "DEFAULT_POSITIVE" || c == "DEFAULT_NEGATIVE"
	}
	return false
}

func isAssignment(s string) bool {
	return len(strings.SplitN(s, "=", 2)) == 2
}

func lineType(s string) LineType {
	if isComment(s) {
		return commentLine
	}

	if isRule(s) {
		return ruleLine
	}

	if isDefaultAssignment(s) {
		return defaultAssignmentLine
	}

	if isAssignment(s) {
		return assignmentLine
	}

	return unknownLine
}
