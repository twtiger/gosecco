package compiler

import (
	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

type JumpsSuite struct{}

var _ = Suite(&JumpsSuite{})

func (s *JumpsSuite) Test_maxSizeJumpSetsUnconditionalJumpPoint(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.BooleanLiteral{true},
			},
			tree.Rule{
				Name: "vhangup",
				Body: tree.BooleanLiteral{true},
			},
			tree.Rule{
				Name: "read",
				Body: tree.BooleanLiteral{true},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t01\t1\n"+
		"jmp\t5\n"+
		"jeq_k\t00\t01\t99\n"+
		"jmp\t3\n"+
		"jeq_k\t00\t01\t0\n"+
		"jmp\t1\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsMulipleUnconditionalJumpPoint(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.BooleanLiteral{true},
			},
			tree.Rule{
				Name: "read",
				Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t01\t1\n"+
		"jmp\t8\n"+
		"jeq_k\t01\t00\t0\n"+
		"jmp\t5\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t01\t02\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsWithTwoComparisons(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
			tree.Rule{
				Name: "read",
				Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t01\t00\t1\n"+
		"jmp\t7\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t01\t00\n"+
		"jmp\tA\n"+
		"jmp\t8\n"+
		"jeq_k\t01\t00\t0\n"+
		"jmp\t5\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t01\t02\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsWithNotEqual(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t01\t00\t1\n"+
		"jmp\t5\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t02\t01\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsWithNotEqualWithMoreThanOneRule(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
			tree.Rule{
				Name: "read",
				Body: tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t01\t00\t1\n"+
		"jmp\t7\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t00\t01\n"+
		"jmp\tA\n"+
		"jmp\t8\n"+
		"jeq_k\t01\t00\t0\n"+
		"jmp\t5\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t02\t01\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}
