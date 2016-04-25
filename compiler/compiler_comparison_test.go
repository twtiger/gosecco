package compiler

import (
	"syscall"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func ComparisonTest(t *testing.T) { TestingT(t) }

type CompilerComparisonSuite struct{}

var _ = Suite(&CompilerComparisonSuite{})

func (s *CompilerComparisonSuite) Test_compilationOfEqualsComparison(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	05	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfSimpleComparisonWithSecondRule(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{Value: 42}},
			},
			tree.Rule{
				Name: "vhangup",
				Body: tree.BooleanLiteral{true},
			},
		},
	}

	res, _ := Compile(p)

	// Load current syscall
	c.Assert(res[0], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    syscallNameIndex,
	})

	// ------------------------- RULE for SYS_WRITE -------------------

	// Compare against the syscall for the current rule
	c.Assert(res[1], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   4,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   2,
		K:    0,
	})

	// Load left hand side of the comparison into A (arg0)
	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	// Compare A against constant K
	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   2,
		Jf:   0,
		K:    42,
	})

	// ------------------------- END RULE for SYS_WRITE -------------------

	// Reload current system call number, since we clobbered A
	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    0,
	})

	// ------------------------- RULE for SYS_VHANGUP -------------------

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   1,
		K:    syscall.SYS_VHANGUP,
	})

	// ------------------------- SHARED RESULT ACTIONS -------------------

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanComparisonToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.GT, Right: tree.NumericLiteral{Value: 42}},
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
		Jf:   5,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   3,
		K:    0,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGT | BPF_K,
		Jt:   0,
		Jf:   1,
		K:    42,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonAToX(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.EQL, Right: tree.Argument{Index: 0}},
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
		Jf:   7,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   5,
		K:    0,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    1,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanComparisonToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.LT, Right: tree.NumericLiteral{Value: 42}},
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
		Jf:   5,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   3,
		K:    0,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGT | BPF_K,
		Jt:   1,
		Jf:   0,
		K:    42,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanOrEqualsToComparisonToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.GTE, Right: tree.NumericLiteral{Value: 42}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGE | BPF_K,
		Jt:   0,
		Jf:   1,
		K:    42,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanOrEqualsToComparisonToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.LTE, Right: tree.NumericLiteral{Value: 42}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGE | BPF_K,
		Jt:   1,
		Jf:   0,
		K:    42,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfNotEqualsToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.NEQL, Right: tree.NumericLiteral{Value: 42}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   1,
		Jf:   0,
		K:    42,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanAToX(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.GT, Right: tree.Argument{Index: 0}},
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
		Jf:   7,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].upper,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   5,
		K:    0,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[0].lower,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    1,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGT | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanOrEqualsToAToX(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.GTE, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGE | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanAToX(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.LT, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGT | BPF_X,
		Jt:   1,
		Jf:   0,
		K:    0,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfLessOrEqualsToAToX(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.LTE, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JGE | BPF_X,
		Jt:   1,
		Jf:   0,
		K:    0,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfNotEqualsAToX(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.NEQL, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP + BPF_JEQ + BPF_X,
		Jt:   1,
		Jf:   0,
		K:    0,
	})
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonForFirstArgument(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{Value: 1}, Op: tree.NEQL, Right: tree.Argument{Index: 1}},
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
		Jf:   7,
		K:    syscall.SYS_WRITE,
	})

	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[1].upper,
	})

	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   5,
		K:    0,
	})

	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    argument[1].lower,
	})

	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    1,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   1,
		Jf:   0,
		K:    0,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[9], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}
