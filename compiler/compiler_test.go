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
