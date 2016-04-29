package compiler

import (
	"testing"

	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func ArithmeticTest(t *testing.T) { TestingT(t) }

type CompilerArithmeticSuite struct{}

var _ = Suite(&CompilerArithmeticSuite{})

func (s *CompilerArithmeticSuite) Test_compilationOfAdditionWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	C\n"+
		"add_k	4\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfMultiplicationWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.MULT,
						Left:  tree.NumericLiteral{3},
						Right: tree.NumericLiteral{8},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	3\n"+
		"mul_k	8\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")

}

func (s *CompilerArithmeticSuite) Test_compilationOfSubtractionWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.MINUS,
						Left:  tree.NumericLiteral{3},
						Right: tree.NumericLiteral{8},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	3\n"+
		"sub_k	8\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfDivisionWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.DIV,
						Left:  tree.NumericLiteral{10},
						Right: tree.NumericLiteral{5},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	A\n"+
		"div_k	5\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfBinaryAndWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.BINAND,
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	4\n"+
		"and_k	2\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfBinaryOrWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.BINOR,
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	4\n"+
		"or_k	2\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfBitwiseLeftShiftWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.LSH,
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	4\n"+
		"lsh_k	2\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfBitwiseRightShiftWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.RSH,
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	4\n"+
		"rsh_k	2\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfModuloWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.MOD,
						Left:  tree.NumericLiteral{10},
						Right: tree.NumericLiteral{3},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	A\n"+
		"mod_k	3\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerArithmeticSuite) Test_compilationOfBinaryXORWithK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left: tree.Argument{Index: 0},
					Op:   tree.EQL,
					Right: tree.Arithmetic{
						Op:    tree.BINXOR,
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{3},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_imm	4\n"+
		"xor_k	3\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}
