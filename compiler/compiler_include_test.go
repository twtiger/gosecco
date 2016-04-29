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
	c.Skip("Need to fix building of jump points to pass this")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.Argument{Index: 0},
					Rights:   []tree.Numeric{tree.NumericLiteral{1}, tree.NumericLiteral{2}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	02	0\n"+ // if it fails, we want to go to test the next element. if it succeeds, we want to test the lower half
		"ld_abs	10\n"+
		"jeq_k	04	00	1\n"+ // jumping directly to final return
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_k	00	01	2\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")

}

func (s *IncludeCompilerSuite) Test_compliationOfNotIncludeOperation(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: false,
					Left:     tree.Argument{Index: 0},
					Rights:   []tree.Numeric{tree.NumericLiteral{1}, tree.NumericLiteral{2}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	06	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	04	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][lower]
		"jeq_k	02	00	1\n"+ // compare to first number in list
		"jeq_k	01	00	2\n"+ // compare to second number in list
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *IncludeCompilerSuite) Test_compliationOfArgumentsInIncludeList(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.NumericLiteral{1},
					Rights:   []tree.Numeric{tree.Argument{Index: 1}, tree.Argument{Index: 0}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+ // syscallNameIndex
		"jeq_k\t00\t0B\t1\n"+ // syscall.SYS_WRITE
		"ld_imm\t1\n"+ // load K into A
		"tax\n"+ // move A to X
		"ld_abs\t1C\n"+ // load first half of argument 1
		"jeq_k\t00\t07\t0\n"+ // compare it to 0
		"ld_abs\t18\n"+ //load second half of argument 1
		"jeq_x\t04\t00\n"+ // compare it to X
		"ld_abs\t14\n"+ // load first half of argument 0
		"jeq_k\t00\t03\t0\n"+ // compare it to 0
		"ld_abs\t10\n"+ // load second half of argument 1
		"jeq_x\t00\t01\n"+ // compare it to X
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *IncludeCompilerSuite) Test_compliationOfIncludeExpressionofNumericWithMixedTypeList(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.NumericLiteral{1},
					Rights:   []tree.Numeric{tree.Argument{Index: 1}, tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+ // syscallNameIndex
		"jeq_k\t00\t09\t1\n"+ // syscall.SYS_WRITE
		"ld_imm\t1\n"+ // load K into A
		"tax\n"+ // move A to X
		"ld_abs\t1C\n"+ // load first half of argument 1
		"jeq_k\t00\t05\t0\n"+ // compare it to 0
		"ld_abs\t18\n"+ //load second half of argument 1
		"jeq_x\t02\t00\n"+ // compare it to X
		"ld_imm\t2A\n"+ // load K into A
		"jeq_x\t00\t01\n"+ // compare it to X
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *IncludeCompilerSuite) Test_compliationOfIncludeExpressionofArgumentWithMixedTypeList(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.Argument{Index: 1},
					Rights:   []tree.Numeric{tree.Argument{Index: 1}, tree.NumericLiteral{0}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	0A	1\n"+ // syscall.SYS_WRITE
		"ld_abs\t1C\n"+ // load argumentindex[0][upper]
		"jeq_k	00	08	0\n"+ // compare to 0
		"ld_abs	18\n"+ // load argumentindex[0][lower]
		"tax\n"+ // move A to X
		"ld_abs\t1C\n"+ // load argumentindex[1][upper]
		"jeq_k\t00\t04\t0\n"+ // compare it to 0
		"ld_abs\t18\n"+ //  load argumentindex[1][lower]
		"jeq_x\t01\t00\n"+ // compare this to X
		"jeq_k\t00\t01\t0\n"+ // compare X against K constant
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}
