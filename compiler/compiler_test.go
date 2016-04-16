package compiler

import (
	"syscall"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

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

	c.Assert(len(res), Equals, 4)

	//   Load 00000000, 00, 00
	c.Assert(res[0], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    0,
	})

	//    Jeq 00000001, 00, 01
	c.Assert(res[1], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   1,
		K:    syscall.SYS_WRITE,
	})

	// Return 7fff0000, 00, 00
	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	// Return 00000000, 00, 00
	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
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

	c.Assert(len(res), Equals, 5)

	//   Load 00000000, 00, 00
	c.Assert(res[0], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    0,
	})

	//    Jeq 00000001, 00, 01
	c.Assert(res[1], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   1,
		Jf:   0,
		K:    syscall.SYS_WRITE,
	})

	//    Jeq 00000099, 00, 01
	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   1,
		K:    syscall.SYS_VHANGUP,
	})

	// Return 7fff0000, 00, 00
	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	// Return 00000000, 00, 00
	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}

func (s *CompilerSuite) Test_compilationOfSimpleComparison(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
			},
			tree.Rule{
				Name: "vhangup",
				Body: tree.BooleanLiteral{true},
			},
		},
	}

	res, _ := Compile(p)

	// spew.Dump(res)

	// c.Assert(len(res), Equals, 5)

	// Load current syscall
	c.Assert(res[0], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    syscallNameIndex,
	})

	// ------------------------- RULE for SYS_WRITE -------------------

	// Compare against the syscall for the current rule
	c.Assert(res[1], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   6,
		Jf:   4,
		K:    syscall.SYS_WRITE,
	})

	// Load the right hand (literal) value of the comparison into A
	c.Assert(res[2], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_IMM,
		K:    42,
	})

	// Move A into the index register to make place for the left hand side in A
	c.Assert(res[3], DeepEquals, unix.SockFilter{
		Code: BPF_MISC | BPF_TAX,
	})

	// Load left hand side of the comparison into A (arg0)
	c.Assert(res[4], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    arg0IndexLowerWord,
	})

	// Compare A against index register (X)
	c.Assert(res[5], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   2,
		Jf:   0,
		K:    0,
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

	// ------------------------- END RULE for SYS_VHANGUP -------------------

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
