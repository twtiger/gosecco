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
	ctx := createCompilerContext()
	compileNumeric(ctx, p)

	c.Assert(asm.Dump(ctx.result), Equals, "ld_imm	2A\n")
}

func (s *NumericCompilerSuite) Test_compilationOfArgument(c *C) {
	ctx := createCompilerContext()
	compileNumeric(ctx, tree.Argument{Type: tree.Low, Index: 3})
	c.Assert(asm.Dump(ctx.result), Equals, "ld_abs	2C\n")

	ctx = createCompilerContext()
	compileNumeric(ctx, tree.Argument{Type: tree.Hi, Index: 1})
	c.Assert(asm.Dump(ctx.result), Equals, "ld_abs	18\n")
}

func (s *NumericCompilerSuite) Test_simpleAdditionOfNumbers(c *C) {
	ctx := createCompilerContext()
	compileNumeric(ctx, tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{3}, Right: tree.NumericLiteral{42}})
	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	2A\n"+
		"st	0\n"+
		"ld_imm	3\n"+
		"ldx_mem	0\n"+
		"add_x\n",
	)
}

// This tests a nested expression:     (((argH1 + 32) * 3) & 42) ^ (argL1 - 15)
func (s *NumericCompilerSuite) Test_moreComplicatedExpression(c *C) {
	ctx := createCompilerContext()
	compileNumeric(ctx,
		tree.Arithmetic{
			Op: tree.BINXOR,
			Left: tree.Arithmetic{
				Op: tree.BINAND,
				Left: tree.Arithmetic{
					Op:    tree.MULT,
					Right: tree.NumericLiteral{3},
					Left: tree.Arithmetic{
						Op:    tree.PLUS,
						Left:  tree.Argument{Type: tree.Hi, Index: 1},
						Right: tree.NumericLiteral{32},
					},
				},
				Right: tree.NumericLiteral{42},
			},
			Right: tree.Arithmetic{
				Op:    tree.MINUS,
				Left:  tree.Argument{Type: tree.Low, Index: 1},
				Right: tree.NumericLiteral{15},
			},
		},
	)
	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	F\n"+
		"st	0\n"+
		"ld_abs	1C\n"+
		"ldx_mem	0\n"+
		"sub_x\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"st	1\n"+
		"ld_imm	3\n"+
		"st	2\n"+
		"ld_imm	20\n"+
		"st	3\n"+
		"ld_abs	18\n"+
		"ldx_mem	3\n"+
		"add_x\n"+
		"ldx_mem	2\n"+
		"mul_x\n"+
		"ldx_mem	1\n"+
		"and_x\n"+
		"ldx_mem	0\n"+
		"xor_x\n",
	)
}
