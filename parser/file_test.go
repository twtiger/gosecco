package parser

import (
	"os"
	"path"
	"strings"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

type FileSuite struct{}

var _ = Suite(&FileSuite{})

func getActualTestFolder() string {
	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, "/parser") {
		return wd
	}
	return path.Join(wd, "parser")
}

func (s *FileSuite) Test_ParseFile(c *C) {
	rp, _ := ParseFile(getActualTestFolder() + "/simple_test_policy")
	c.Assert(rp, DeepEquals, tree.RawPolicy{
		RuleOrMacros: []interface{}{
			tree.Macro{
				Name:          "DEFAULT_POSITIVE",
				ArgumentNames: nil,
				Body:          tree.Variable{Name: "kill"}},
			tree.Macro{
				Name:          "something",
				ArgumentNames: []string{"a"},
				Body:          tree.Arithmetic{Op: 0, Left: tree.NumericLiteral{Value: 0x1}, Right: tree.Variable{Name: "a"}}},
			tree.Macro{
				Name:          "VAL",
				ArgumentNames: nil,
				Body:          tree.NumericLiteral{Value: 0x2a}},
			tree.Rule{
				Name:           "read",
				PositiveAction: "",
				NegativeAction: "",
				Body:           tree.NumericLiteral{Value: 0x2a}},
		}})
}

func (s *FileSuite) Test_ParseFile_failing(c *C) {
	rp, ee := ParseFile(getActualTestFolder() + "/failing_test_policy")
	c.Assert(rp.RuleOrMacros, IsNil)
	c.Assert(ee, ErrorMatches, ".*parser/failing_test_policy:1: unexpected end of line")
}
