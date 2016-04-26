package unifier

import (
	"testing"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func ActionsTest(t *testing.T) { TestingT(t) }

type UnifierActionsSuite struct{}

var _ = Suite(&UnifierActionsSuite{})

func (s *UnifierActionsSuite) Test_Unify_setsDefaultActionsForEnforcedWhitelist(c *C) {
	input := tree.RawPolicy{
		ListType:     tree.WhiteList,
		RuleOrMacros: []interface{}{},
	}

	output, _ := Unify(input, true)

	c.Assert(output.DefaultPositiveAction, Equals, "allow")
	c.Assert(output.DefaultNegativeAction, Equals, "kill")
}

func (s *UnifierActionsSuite) Test_Unify_setsDefaultActionsForEnforcedBlacklist(c *C) {
	input := tree.RawPolicy{
		ListType:     tree.BlackList,
		RuleOrMacros: []interface{}{},
	}

	output, _ := Unify(input, true)

	c.Assert(output.DefaultPositiveAction, Equals, "kill")
	c.Assert(output.DefaultNegativeAction, Equals, "allow")
}

func (s *UnifierActionsSuite) Test_Unify_setsDefaultActionsForNotEnforcedWhitelist(c *C) {
	input := tree.RawPolicy{
		ListType:     tree.WhiteList,
		RuleOrMacros: []interface{}{},
	}

	output, _ := Unify(input, false)

	c.Assert(output.DefaultPositiveAction, Equals, "allow")
	c.Assert(output.DefaultNegativeAction, Equals, "trace")
}

func (s *UnifierActionsSuite) Test_Unify_setsDefaultActionsForNotEnforcedBlacklist(c *C) {
	input := tree.RawPolicy{
		ListType:     tree.BlackList,
		RuleOrMacros: []interface{}{},
	}

	output, _ := Unify(input, false)

	c.Assert(output.DefaultPositiveAction, Equals, "trace")
	c.Assert(output.DefaultNegativeAction, Equals, "allow")
}
