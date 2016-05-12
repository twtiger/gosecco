package compiler2

import (
	"syscall"

	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

type CompilerSuite struct{}

var _ = Suite(&CompilerSuite{})

func (s *CompilerSuite) Test_simplestCompilation(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.BooleanLiteral{true},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	01	00	1\n"+
		"ret_k	0\n"+
		"ret_k	7FFF0000\n")
}

func (s *CompilerSuite) Test_nextSimplestCompilation(c *C) {
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
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	02	00	1\n"+
		"jeq_k	01	00	99\n"+
		"ret_k	0\n"+
		"ret_k	7FFF0000\n")
}

func (s *CompilerSuite) Test_stackOverflowCreatesError(c *C) {
	cx := createCompilerContext()
	cx.stackTop = syscall.BPF_MEMWORDS
	c.Assert(cx.pushAToStack(), ErrorMatches, "the expression is too complicated to compile. Please refer to the language documentation")
}

func (s *CompilerSuite) Test_stackDoesNotOverflowRightBeforeItsLimit(c *C) {
	cx := createCompilerContext()
	cx.stackTop = syscall.BPF_MEMWORDS - 1
	c.Assert(cx.pushAToStack(), IsNil)
}

func (s *CompilerSuite) Test_stackDoesNotPopAfterReachingTheLowestIndex(c *C) {
	cx := createCompilerContext()
	cx.stackTop = 0
	c.Assert(cx.popStackToX(), ErrorMatches, "popping from empty stack - this is likely a programmer error")
}

func (s *CompilerSuite) Test_maxSizeJumpSetsUnconditionalJumpPoint(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := []tree.Rule{
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
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k\t00\t01\t1\n"+
		"jmp\t3\n"+
		"jeq_k\t02\t00\t99\n"+
		"jeq_k\t01\t00\t0\n"+
		"ret_k\t0\n"+
		"ret_k\t7FFF0000\n")
}

func (s *CompilerSuite) Test_maxSizeJumpSetsMulipleUnconditionalJumpPoint(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := []tree.Rule{
		tree.Rule{
			Name: "write",
			Body: tree.BooleanLiteral{true},
		},
		tree.Rule{
			Name: "read",
			Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
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
		"jeq_x\t01\t00\n"+
		"ret_k\t0\n"+
		"ret_k\t7FFF0000\n")
}

func (s *CompilerSuite) Test_maxSizeJumpSetsWithTwoComparisons(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := []tree.Rule{
		tree.Rule{
			Name: "write",
			Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
		},
		tree.Rule{
			Name: "read",
			Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t01\t00\t1\n"+
		"jmp\t6\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t00\t01\n"+
		"jmp\t8\n"+
		"jeq_k\t01\t00\t0\n"+
		"jmp\t5\n"+
		"ld_imm\t1\n"+
		"st\t0\n"+
		"ld_imm\t2A\n"+
		"ldx_mem\t0\n"+
		"jeq_x\t01\t00\n"+
		"ret_k\t0\n"+
		"ret_k\t7FFF0000\n")
}
