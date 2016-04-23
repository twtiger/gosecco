package compiler

import (
	"testing"

	"github.com/twtiger/gosecco/asm"

	"github.com/twtiger/gosecco/tree"

	. "gopkg.in/check.v1"
)

func IncludeTest(t *testing.T) { TestingT(t) }

type IncludeCompilerSuite struct{}

var _ = Suite(&IncludeCompilerSuite{})

func (s *IncludeCompilerSuite) Test_compliationOfIncludeOperation(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Left:   tree.Argument{Index: 0},
					Rights: []tree.Numeric{tree.NumericLiteral{1}, tree.NumericLiteral{2}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	06	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	04	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][lower]
		"jeq_k	01	00	1\n"+ // compare to first number in list
		"jeq_k	00	01	2\n"+ // compare to second number in list
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}
