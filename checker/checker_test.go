package checker

import (
	"testing"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CheckerSuite struct{}

var _ = Suite(&CheckerSuite{})

func (s *CheckerSuite) Test_checksNumber(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.NumericLiteral{Value: 42}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: 42")
}

func (s *CheckerSuite) Test_checksComparisonRight(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{Value: 42}, Right: tree.BooleanLiteral{false}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_checksComparisonLeft(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}, Left: tree.BooleanLiteral{false}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_checksSuccessfulSimpleCase(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.BooleanLiteral{true}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 0)
}

func (s *CheckerSuite) Test_checksComparisonInNumericContext(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL,
			Left:  tree.NumericLiteral{Value: 23},
			Right: tree.Comparison{Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}, Left: tree.NumericLiteral{Value: 15}}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: \\(eq 15 42\\)")
}

func (s *CheckerSuite) Test_checksAndAsArgument(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL,
			Left:  tree.NumericLiteral{Value: 23},
			Right: tree.And{Left: tree.BooleanLiteral{true}, Right: tree.BooleanLiteral{true}}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: \\(and true true\\)")
}

func (s *CheckerSuite) Test_checksAndArgumentLeft(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.And{
			Left:  tree.NumericLiteral{Value: 23},
			Right: tree.BooleanLiteral{true},
		}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: 23")
}

func (s *CheckerSuite) Test_checksAndArgumentRight(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.And{
			Right: tree.NumericLiteral{Value: 23},
			Left:  tree.BooleanLiteral{true},
		}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: 23")
}

func (s *CheckerSuite) Test_checksOrAsArgument(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL,
			Left:  tree.NumericLiteral{Value: 23},
			Right: tree.Or{Left: tree.BooleanLiteral{true}, Right: tree.BooleanLiteral{true}}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: \\(or true true\\)")
}

func (s *CheckerSuite) Test_checksOrArgumentLeft(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Or{
			Left:  tree.NumericLiteral{Value: 23},
			Right: tree.BooleanLiteral{true},
		}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: 23")
}

func (s *CheckerSuite) Test_checksOrArgumentRight(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Or{
			Right: tree.NumericLiteral{Value: 23},
			Left:  tree.BooleanLiteral{true},
		}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: 23")
}

func (s *CheckerSuite) Test_argumentShouldTypecheckAsNumeric(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Argument{Index: 0, Type: tree.Hi}}}}
	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: argH0")
}

func (s *CheckerSuite) Test_checksSuccessfulSimpleCaseWithArg(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}, Left: tree.Argument{Index: 0}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 0)
}

func (s *CheckerSuite) Test_checksInvalidVariable(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Variable{"foo"}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] found unresolved variable: foo")
}

func (s *CheckerSuite) Test_checksInvalidCall(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Call{Name: "foox", Args: nil}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] found unresolved call: foox")
}

func (s *CheckerSuite) Test_checksSuccessfulSimpleCaseWithBinNeg(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}, Left: tree.BinaryNegation{tree.NumericLiteral{Value: 42}}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 0)
}

func (s *CheckerSuite) Test_checksInvalidSimpleArgumentToBinNeg(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}, Left: tree.BinaryNegation{tree.BooleanLiteral{false}}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_checksInvalidCaseWithBinNeg(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.BinaryNegation{tree.NumericLiteral{Value: 42}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: \\(binNeg 42\\)")
}

func (s *CheckerSuite) Test_checksArithRight(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Left: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{Value: 42}, Right: tree.BooleanLiteral{false}}, Right: tree.NumericLiteral{Value: 42}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_checksArithLeft(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Left: tree.Arithmetic{Op: tree.PLUS, Right: tree.NumericLiteral{Value: 42}, Left: tree.BooleanLiteral{false}}, Right: tree.NumericLiteral{Value: 5233}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_checksArithAsArgument(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Arithmetic{Op: tree.PLUS,
			Left:  tree.NumericLiteral{Value: 23},
			Right: tree.NumericLiteral{Value: 24},
		}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: \\(plus 23 24\\)")
}

func (s *CheckerSuite) Test_checksSuccess_booleanNegation(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Negation{tree.BooleanLiteral{false}}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 0)
}

func (s *CheckerSuite) Test_checksInvalidPlacement_booleanNegation(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}, Left: tree.Negation{tree.BooleanLiteral{false}}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: \\(not false\\)")
}

func (s *CheckerSuite) Test_checksInvalidArgument_booleanNegation(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Negation{tree.NumericLiteral{Value: 42}}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected boolean expression but found: 42")
}

func (s *CheckerSuite) Test_inclusion_success(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Inclusion{Positive: true, Left: tree.NumericLiteral{Value: 42}, Rights: []tree.Numeric{tree.NumericLiteral{Value: 42}, tree.NumericLiteral{Value: 42}, tree.NumericLiteral{Value: 42}}}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 0)
}

func (s *CheckerSuite) Test_inclusion_badPlacement(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Comparison{Op: tree.EQL, Left: tree.Inclusion{Positive: true, Left: tree.NumericLiteral{Value: 42}, Rights: []tree.Numeric{tree.NumericLiteral{Value: 42}, tree.NumericLiteral{Value: 42}, tree.NumericLiteral{Value: 42}}}, Right: tree.NumericLiteral{Value: 42}}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: \\(in 42 42 42 42\\)")
}

func (s *CheckerSuite) Test_inclusion_badLeft(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Inclusion{Positive: true, Left: tree.BooleanLiteral{false}, Rights: []tree.Numeric{tree.NumericLiteral{Value: 42}, tree.NumericLiteral{Value: 42}, tree.NumericLiteral{Value: 42}}}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_inclusion_badRight(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.Inclusion{Positive: true, Left: tree.NumericLiteral{Value: 23}, Rights: []tree.Numeric{tree.NumericLiteral{Value: 42}, tree.BooleanLiteral{false}, tree.NumericLiteral{Value: 42}}}}}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[read\\] expected numeric expression but found: false")
}

func (s *CheckerSuite) Test_duplicateRules(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "read", Body: tree.BooleanLiteral{true}},
		tree.Rule{Name: "write", Body: tree.BooleanLiteral{true}},
		tree.Rule{Name: "read", Body: tree.BooleanLiteral{true}},
		tree.Rule{Name: "write", Body: tree.BooleanLiteral{true}},
		tree.Rule{Name: "fcntl", Body: tree.BooleanLiteral{true}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 2)
	c.Assert(val[0], ErrorMatches, "\\[read\\] duplicate definition of syscall rule")
	c.Assert(val[1], ErrorMatches, "\\[write\\] duplicate definition of syscall rule")
}

func (s *CheckerSuite) Test_invalidSyscall(c *C) {
	toCheck := tree.Policy{Rules: []tree.Rule{
		tree.Rule{Name: "fluffipuff", Body: tree.BooleanLiteral{true}},
	}}

	val := EnsureValid(toCheck)

	c.Assert(len(val), Equals, 1)
	c.Assert(val[0], ErrorMatches, "\\[fluffipuff\\] invalid syscall")
}
