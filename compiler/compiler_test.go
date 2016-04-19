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

func (s *CompilerSuite) Test_compliationOfArtimeticOperation(c *C) {
	c.Skip("finish implementation to make this pass")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
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
		Jt:   6,
		Jf:   7,
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
		K:    arg0IndexLowerWord,
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_X,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_ALLOW,
	})

	c.Assert(res[8], DeepEquals, unix.SockFilter{
		Code: BPF_RET | BPF_K,
		K:    SECCOMP_RET_KILL,
	})
}
