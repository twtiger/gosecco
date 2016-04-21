package compiler

import (
	"syscall"
	"testing"

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
					Left:  tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
					Right: tree.Comparison{Left: tree.Argument{0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
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
		K:    ArgumentIndex[0]["upper"],
	})

	c.Assert(res[6], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   10,
		K:    0,
	})

	c.Assert(res[7], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    ArgumentIndex[0]["lower"],
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
		K:    ArgumentIndex[0]["upper"],
	})

	c.Assert(res[13], DeepEquals, unix.SockFilter{
		Code: BPF_JMP | BPF_JEQ | BPF_K,
		Jt:   0,
		Jf:   3,
		K:    0,
	})

	c.Assert(res[14], DeepEquals, unix.SockFilter{
		Code: BPF_LD | BPF_W | BPF_ABS,
		K:    ArgumentIndex[0]["lower"],
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
