package simplify

import (
	"testing"

	"github.com/twtiger/go-seccomp/tree"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SimplifierSuite struct{}

var _ = Suite(&SimplifierSuite{})

func (s *SimplifierSuite) Test_simplifyAddition(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{1}, Right: tree.NumericLiteral{2}})
	c.Assert(tree.ExpressionString(sx), Equals, "3")
}

func (s *SimplifierSuite) Test_simplifySubtraction(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.MINUS, Left: tree.NumericLiteral{32}, Right: tree.NumericLiteral{3}})
	c.Assert(tree.ExpressionString(sx), Equals, "29")
}

func (s *SimplifierSuite) Test_simplifyMult(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.MULT, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{3}})
	c.Assert(tree.ExpressionString(sx), Equals, "36")
}

func (s *SimplifierSuite) Test_simplifyDiv(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.DIV, Left: tree.NumericLiteral{37}, Right: tree.NumericLiteral{3}})
	c.Assert(tree.ExpressionString(sx), Equals, "12")
}

func (s *SimplifierSuite) Test_simplifyMod(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.MOD, Left: tree.NumericLiteral{37}, Right: tree.NumericLiteral{3}})
	c.Assert(tree.ExpressionString(sx), Equals, "1")
}

func (s *SimplifierSuite) Test_simplifyBinAnd(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.BINAND, Left: tree.NumericLiteral{7}, Right: tree.NumericLiteral{4}})
	c.Assert(tree.ExpressionString(sx), Equals, "4")
}

func (s *SimplifierSuite) Test_simplifyBinOr(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.BINOR, Left: tree.NumericLiteral{3}, Right: tree.NumericLiteral{8}})
	c.Assert(tree.ExpressionString(sx), Equals, "11")
}

func (s *SimplifierSuite) Test_simplifyBinXoe(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.BINXOR, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{12}})
	c.Assert(tree.ExpressionString(sx), Equals, "38")
}

func (s *SimplifierSuite) Test_simplifyLsh(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.LSH, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{2}})
	c.Assert(tree.ExpressionString(sx), Equals, "168")
}

func (s *SimplifierSuite) Test_simplifyRsh(c *C) {
	sx := Simplify(tree.Arithmetic{Op: tree.RSH, Left: tree.NumericLiteral{84}, Right: tree.NumericLiteral{2}})
	c.Assert(tree.ExpressionString(sx), Equals, "21")
}
