package compiler

import (
	"syscall"
	"testing"

	"github.com/twtiger/gosecco/asm"
	"golang.org/x/sys/unix"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func BoolTest(t *testing.T) { TestingT(t) }

type BoolCompilerSuite struct{}

var _ = Suite(&BoolCompilerSuite{})

func (s *BoolCompilerSuite) Test_compliationOfOrOperation(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Or{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
					Right: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
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
		Jf:   15,
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
		Jf:   10,
		K:    0,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   7,
		Jf:   0,
		K:    0,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    12,
	})

	c.Assert(res[10], DeepEquals, unix.SockFilter{
		Code: BPF_ALU | BPF_ADD | BPF_K,
		K:    4,
	})

	c.Assert(res[11], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	c.Assert(res[12], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[13], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   3,
		K:    0,
	})

	c.Assert(res[14], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[15], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	c.Assert(res[16], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[17], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfOrExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Or{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	04	00	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfAndExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	00	05	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfNegatedExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Negation{
					Operand: tree.And{
						Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
						Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	05	00	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfNestedNegatedExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Negation{
						Operand: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	00	05	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}
