package compiler2

import (
	"syscall"

	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

type BooleanCompilerSuite struct{}

var _ = Suite(&BooleanCompilerSuite{})

func (s *BooleanCompilerSuite) Test_compilationOfSimpleComparison(c *C) {
	p := tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "pos", "neg")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"pos": []int{4},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"neg": []int{4},
	})
}

func (s *BooleanCompilerSuite) Test_compilationOfSimpleComparison2(c *C) {
	p := tree.Comparison{Op: tree.NEQL, Left: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{25}}, Right: tree.NumericLiteral{1}}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "posx", "negx")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	19\n"+
		"st	1\n"+
		"ld_imm	2A\n"+
		"ldx_mem	1\n"+
		"add_x\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"negx": []int{8},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"posx": []int{8},
	})
}

func (s *BooleanCompilerSuite) Test_compilationOfSimpleComparison3(c *C) {
	p := tree.Comparison{Op: tree.GT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "pos", "neg")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jgt_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"pos": []int{4},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"neg": []int{4},
	})
}

func (s *BooleanCompilerSuite) Test_compilationOfSimpleComparison4(c *C) {
	p := tree.Comparison{Op: tree.GTE, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "pos", "neg")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jge_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"pos": []int{4},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"neg": []int{4},
	})
}

func (s *BooleanCompilerSuite) Test_compilationOfInvalidComparison(c *C) {
	p := tree.Comparison{Op: tree.LT, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}}
	ctx := createCompilerContext()
	res := compileBoolean(ctx, p, false, "pos", "neg")
	c.Assert(res, Not(IsNil))
	c.Assert(res, ErrorMatches, "this comparison type is not allowed - this is probably a programmer error: \\(lt 42 1\\)")
}

func (s *BooleanCompilerSuite) Test_compilationOfSimpleAnd(c *C) {
	p := tree.And{
		Left:  tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
		Right: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{41}, Right: tree.NumericLiteral{23}},
	}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "pos", "neg")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n"+
		"ld_imm	17\n"+
		"st	0\n"+
		"ld_imm	29\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"pos":               []int{9},
		"generatedLabel000": []int{4},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"neg": []int{4, 9},
	})
	c.Assert(ctx.labels, DeepEquals, map[label]int{
		"generatedLabel000": 5,
	})
}

func (s *BooleanCompilerSuite) Test_compilationOfSimpleOr(c *C) {
	p := tree.Or{
		Left:  tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
		Right: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{41}, Right: tree.NumericLiteral{23}},
	}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "pos", "neg")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n"+
		"ld_imm	17\n"+
		"st	0\n"+
		"ld_imm	29\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"pos": []int{4, 9},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"generatedLabel000": []int{4},
		"neg":               []int{9},
	})
	c.Assert(ctx.labels, DeepEquals, map[label]int{
		"generatedLabel000": 5,
	})
}

func (s *BooleanCompilerSuite) Test_compilationOfSimpleNegation(c *C) {
	p := tree.Negation{
		Operand: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
	}
	ctx := createCompilerContext()
	compileBoolean(ctx, p, false, "pos", "neg")

	c.Assert(asm.Dump(ctx.result), Equals, ""+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	00\n",
	)
	c.Assert(ctx.jts, DeepEquals, map[label][]int{
		"neg": []int{4},
	})
	c.Assert(ctx.jfs, DeepEquals, map[label][]int{
		"pos": []int{4},
	})
	c.Assert(ctx.labels, DeepEquals, map[label]int{})
}

func (s *BooleanCompilerSuite) Test_thatAnErrorIsSetWhenWeCompileAfterReachingTheMaximumHeightOfTheStack(c *C) {
	ctx := createCompilerContext()
	ctx.stackTop = syscall.BPF_MEMWORDS
	p := tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}}
	err := compileBoolean(ctx, p, false, "pos", "neg")
	c.Assert(err, ErrorMatches, "the expression is too complicated to compile. Please refer to the language documentation.")
}
