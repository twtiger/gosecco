package compiler

import (
	"syscall"
	"testing"

	"golang.org/x/sys/unix"

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

	c.Assert(res[0], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    syscallNameIndex,
	})

	c.Assert(res[1], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   8,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    12,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_ALU | BPF_ADD | BPF_K,
		K:    4,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   3,
		K:    0,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[10], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
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

	c.Assert(res[0], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    syscallNameIndex,
	})

	c.Assert(res[1], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   8,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    3,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_ALU | BPF_MUL | BPF_K,
		K:    8,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   3,
		K:    0,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[10], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
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
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_imm\t3\n"+
		"sub_k\t8\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t10\n"+
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
						Left:  tree.NumericLiteral{10},
						Right: tree.NumericLiteral{5},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_imm\tA\n"+
		"div_k\t5\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t10\n"+
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
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_imm\t4\n"+
		"and_k\t2\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t10\n"+
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
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_imm\t4\n"+
		"or_k\t2\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t10\n"+
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
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_imm\t4\n"+
		"lsh_k\t2\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t10\n"+
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
						Left:  tree.NumericLiteral{4},
						Right: tree.NumericLiteral{2},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t08\t1\n"+
		"ld_imm\t4\n"+
		"rsh_k\t2\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t01\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}
