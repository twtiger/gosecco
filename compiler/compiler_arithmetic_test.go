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
				Body: tree.Comparison{Left: tree.Argument{Index: 0},
					Op:    tree.EQL,
					Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.Argument{Index: 1, Type: tree.Low}, Right: tree.NumericLiteral{4}}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"add_k\t4\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{8},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"mul_k\t8\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{8},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"sub_k\t8\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{5},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"div_k\t5\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"and_k\t2\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"or_k\t2\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"lsh_k\t2\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"rsh_k\t2\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{3},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"mod_k\t3\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
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
						Left:  tree.Argument{Index: 1, Type: tree.Low},
						Right: tree.NumericLiteral{3},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_abs\t1C\n"+
		"xor_k\t3\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t03\n"+
		"ld_abs\t14\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}
