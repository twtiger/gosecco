package compiler2

import (
	"testing"

	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type NumericCompilerSuite struct{}

var _ = Suite(&NumericCompilerSuite{})

func (s *NumericCompilerSuite) Test_compilationOfLiteral(c *C) {
	p := tree.NumericLiteral{42}
	ctx := &compilerContext{}
	compileNumeric(ctx, p)

	c.Assert(asm.Dump(ctx.result), Equals, "ld_imm	2A\n")
}
