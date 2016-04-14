package parser

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/twtiger/go-seccomp/tree"
)

func parseLines(lines []string) (tree.RawPolicy, error) {
	// TODO: keep track of line numbers for errors

	result := []interface{}{}

	for _, l := range lines {
		switch lineType(l) {
		case commentLine: //ignore
		case emptyLine: //ignore
		case ruleLine:
			parsedRule, err := parseRule(l)
			if err != nil {
				return tree.RawPolicy{}, err
			}
			result = append(result, parsedRule)
		case assignmentLine, defaultAssignmentLine:
			// TODO: parse assignment
		case unknownLine:
			return tree.RawPolicy{}, fmt.Errorf("Couldn't parse line: '%s' - it doesn't match any kind of valid syntax", l)
		}
	}

	return tree.RawPolicy{result}, nil
}

func parseFile(path string) (tree.RawPolicy, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return tree.RawPolicy{}, err
	}
	return parseLines(strings.Split(string(file), "\n"))
}
