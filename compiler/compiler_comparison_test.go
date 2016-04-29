package compiler

import (
	"testing"

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
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
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
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
			},
			tree.Rule{
				Name: "vhangup",
				Body: tree.BooleanLiteral{true},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	04	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	10\n"+
		"jeq_k	02	00	2A\n"+
		"ld_abs	0\n"+
		"jeq_k	00	01	99\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanComparisonToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.GT, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jgt_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonAToX(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.EQL, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	10\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanComparisonToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.LT, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jgt_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanOrEqualsToComparisonToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.GTE, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jge_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanOrEqualsToComparisonToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.LTE, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jge_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfNotEqualsToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.NEQL, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanAToX(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.GT, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	10\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jgt_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanOrEqualsToAToX(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.GTE, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	10\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jge_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanAToX(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.LT, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	10\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jgt_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfLessOrEqualsToAToX(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.LTE, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	10\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jge_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfNotEqualsAToX(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.NEQL, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	10\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jeq_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonForFirstArgument(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.NEQL, Right: tree.Argument{Index: 1}},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	1C\n"+
		"jeq_k	00	05	0\n"+
		"ld_abs	18\n"+
		"tax\n"+
		"ld_imm	1\n"+
		"jeq_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}
