package unifier

import (
	"testing"

	"github.com/twtiger/go-seccomp/tree"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type UnifierSuite struct{}

var _ = Suite(&UnifierSuite{})

func (s *UnifierSuite) Test_Unify_withNothingToUnify(c *C) {
	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{},
	}

	output := Unify(input)

	c.Assert(output.DefaultPositiveAction, Equals, "")
	c.Assert(output.DefaultNegativeAction, Equals, "")
	c.Assert(len(output.Macros), Equals, 0)
	c.Assert(len(output.Rules), Equals, 0)
}

func (s *UnifierSuite) Test_Unify_withRuleToUnify(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			rule,
		},
	}

	output := Unify(input)

	c.Assert(output.DefaultPositiveAction, Equals, "")
	c.Assert(output.DefaultNegativeAction, Equals, "")
	c.Assert(len(output.Macros), Equals, 0)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(output.Rules[0], Equals, rule)
}

func (s *UnifierSuite) Test_Unify_withRuleAndMacroThatDoesntUnify(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
	}

	macro := tree.Macro{
		Name: "var1",
		Body: tree.NumericLiteral{1},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			rule,
			macro,
		},
	}

	output := Unify(input)

	c.Assert(output.DefaultPositiveAction, Equals, "")
	c.Assert(output.DefaultNegativeAction, Equals, "")
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["var1"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(output.Rules[0], DeepEquals, rule)
}

func (s *UnifierSuite) Test_Unify_withRuleAndMacroToActuallyUnify(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.Variable{"var1"}},
	}

	macro := tree.Macro{
		Name: "var1",
		Body: tree.NumericLiteral{1},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["var1"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(eq arg0 1)")
}

func (s *UnifierSuite) Test_Unify_orExpression(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Or{Left: tree.Argument{0}, Right: tree.Variable{"var1"}},
	}

	macro1 := tree.Macro{
		Name: "var1",
		Body: tree.NumericLiteral{1},
	}

	macro2 := tree.Macro{
		Name: "var2",
		Body: tree.NumericLiteral{2},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro1,
			macro2,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 2)
	c.Assert(output.Macros["var1"], DeepEquals, macro1)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(or arg0 1)")
}

func (s *UnifierSuite) Test_Unify_withAndExpressione(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.And{Left: tree.Argument{0}, Right: tree.Variable{"var1"}},
	}

	macro := tree.Macro{
		Name: "var1",
		Body: tree.NumericLiteral{1},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["var1"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(and arg0 1)")
}

func (s *UnifierSuite) Test_Unify_withArithmeticExpression(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Arithmetic{Left: tree.Argument{0}, Op: tree.PLUS, Right: tree.Variable{"var1"}},
	}

	macro := tree.Macro{
		Name: "var1",
		Body: tree.NumericLiteral{1},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["var1"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(plus arg0 1)")
}

func (s *UnifierSuite) Test_Unify_withInclusionExpression(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Inclusion{Positive: true,
			Left:   tree.Argument{0},
			Rights: []tree.Numeric{tree.NumericLiteral{1}, tree.Variable{"var2"}},
		},
	}

	macro := tree.Macro{
		Name: "var2",
		Body: tree.NumericLiteral{2},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["var2"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(in arg0 1 2)")
}

func (s *UnifierSuite) Test_Unify_withInclusionExpressionVariableLeft(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Inclusion{Positive: true,
			Left:   tree.Variable{"var1"},
			Rights: []tree.Numeric{tree.NumericLiteral{1}, tree.Variable{"var2"}},
		},
	}

	macro1 := tree.Macro{
		Name: "var1",
		Body: tree.Argument{0},
	}

	macro2 := tree.Macro{
		Name: "var2",
		Body: tree.NumericLiteral{2},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro1,
			macro2,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 2)
	c.Assert(output.Macros["var1"], DeepEquals, macro1)
	c.Assert(output.Macros["var2"], DeepEquals, macro2)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(in arg0 1 2)")
}

func (s *UnifierSuite) Test_Unify_withNegationExpression(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Negation{Operand: tree.Variable{"var1"}},
	}

	macro := tree.Macro{
		Name: "var1",
		Body: tree.BooleanLiteral{true},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["var1"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(not true)")
}

func (s *UnifierSuite) Test_Unify_withCallExpression(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Call{Name: "compV1", Args: []tree.Any{tree.Argument{0}}},
	}

	macro := tree.Macro{
		Name:          "compV1",
		ArgumentNames: []string{"var1"},
		Body:          tree.Comparison{Left: tree.Variable{"var1"}, Op: tree.EQL, Right: tree.NumericLiteral{1}},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	//c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["compV1"], DeepEquals, macro)
	//c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(eq arg0 1)")
}

func (s *UnifierSuite) Test_Unify_withCallExpressionWithMultipleVariables(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Call{Name: "compV1", Args: []tree.Any{tree.Argument{0}, tree.Argument{1}}},
	}

	macro := tree.Macro{
		Name:          "compV1",
		ArgumentNames: []string{"var1", "var2"},
		Body:          tree.Comparison{Left: tree.Variable{"var1"}, Op: tree.EQL, Right: tree.Variable{"var2"}},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 1)
	c.Assert(output.Macros["compV1"], DeepEquals, macro)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(eq arg0 arg1)")
}

func (s *UnifierSuite) Test_Unify_withCallExpressionWithPreviouslyDefinedVariables(c *C) {
	rule := tree.Rule{
		Name: "write",
		Body: tree.Call{Name: "compV1", Args: []tree.Any{tree.Argument{0}, tree.Variable{"var1"}}},
	}

	macro1 := tree.Macro{
		Name: "var1",
		Body: tree.Argument{5},
	}

	macro2 := tree.Macro{
		Name:          "compV1",
		ArgumentNames: []string{"var1", "var2"},
		Body:          tree.Comparison{Left: tree.Variable{"var1"}, Op: tree.EQL, Right: tree.Variable{"var2"}},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			macro1,
			macro2,
			rule,
		},
	}

	output := Unify(input)
	c.Assert(output.Macros["var1"], DeepEquals, macro1)
	c.Assert(output.Macros["compV1"], DeepEquals, macro2)
	c.Assert(len(output.Rules), Equals, 1)
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(eq arg0 arg5)")
}

func (s *UnifierSuite) Test_Unify_withNoVariableDefinedRaisesNoVariableDefinedError(c *C) {
	c.Skip("handle panic")
	rule := tree.Rule{
		Name: "write",
		Body: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.Variable{"var1"}},
	}

	input := tree.RawPolicy{
		RuleOrMacros: []interface{}{
			rule,
		},
	}

	output := Unify(input)
	c.Assert(len(output.Macros), Equals, 0)
	c.Assert(len(output.Rules), Equals, 0)
	// TODO fix this up
	c.Assert(tree.ExpressionString(output.Rules[0].Body), Equals, "(eq arg0 1)")
}
