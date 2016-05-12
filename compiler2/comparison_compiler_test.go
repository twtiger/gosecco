package compiler2

import (
	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

type ComparisonCompilerSuite struct{}

var _ = Suite(&ComparisonCompilerSuite{})

func (s *ComparisonCompilerSuite) Test_SingleComparisons(c *C) {
	ctx := createCompilerContext()

	p := []tree.Rule{
		tree.Rule{
			Name: "write",
			Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t05\t1\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t01\t02\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}
