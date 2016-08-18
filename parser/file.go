package parser

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/twtiger/gosecco/tree"
)

// ParseError represents error parsing a policy file. It will report the filename and the line number as well as the actual error.
type ParseError struct {
	originalError error
	file          string
	line          int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.file, e.line, e.originalError)
}

func parseLines(path string, lines []string) (tree.RawPolicy, error) {
	result := []interface{}{}

	for ix, l := range lines {
		switch lineType(l) {
		case commentLine: //ignore
		case emptyLine: //ignore
		case ruleLine:
			parsedRule, err := parseRule(l)
			if err != nil {
				return tree.RawPolicy{}, &ParseError{err, path, ix}
			}
			result = append(result, parsedRule)
		case assignmentLine, defaultAssignmentLine:
			parsedBinding, err := parseBinding(l)
			if err != nil {
				return tree.RawPolicy{}, &ParseError{err, path, ix}
			}
			result = append(result, parsedBinding)

		case unknownLine:
			return tree.RawPolicy{}, &ParseError{fmt.Errorf("Couldn't parse line: '%s' - it doesn't match any kind of valid syntax", l), path, ix}
		}
	}

	return tree.RawPolicy{RuleOrMacros: result}, nil
}

// ParseFile will parse the given file and return a raw parse tree or the error generated
func ParseFile(path string) (tree.RawPolicy, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return tree.RawPolicy{}, err
	}
	return parseLines(path, strings.Split(string(file), "\n"))
}

// ParseString will parse the given string and return a raw parse tree or the error generated
func ParseString(str string) (tree.RawPolicy, error) {
	return parseLines("<string>", strings.Split(string(str), "\n"))
}
