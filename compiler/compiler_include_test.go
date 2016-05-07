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
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	04	00	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jeq_k	00	01	2\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *IncludeCompilerSuite) Test_compliationOfNotIncludeOperation(c *C) {
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
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	05	00	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	01	00	2\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *IncludeCompilerSuite) Test_compliationOfArgumentsInIncludeList(c *C) {
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
		"ld_abs\t0\n"+
		"jeq_k\t00\t09\t1\n"+
		"ld_abs\t18\n"+
		"jeq_k\t00\t02\t0\n"+
		"ld_abs\t1C\n"+
		"jeq_k\t04\t00\t1\n"+
		"ld_abs\t10\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t14\n"+
		"jeq_k\t00\t01\t1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *IncludeCompilerSuite) Test_compliationOfIncludeExpressionofArgumentWithMixedTypeList(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.Argument{Index: 1},
					Rights:   []tree.Numeric{tree.Argument{Index: 0}, tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t0D\t1\n"+
		"ld_abs\t18\n"+
		"tax\n"+
		"ld_abs\t10\n"+
		"jeq_x\t00\t04\n"+
		"ld_abs\t1C\n"+
		"tax\n"+
		"ld_abs\t14\n"+
		"jeq_x\t04\t00\n"+
		"ld_abs\t18\n"+
		"jeq_k\t00\t03\t0\n"+
		"ld_abs\t1C\n"+
		"jeq_k\t00\t01\t2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *IncludeCompilerSuite) Test_compliationOfIncludeExpressionofExpressionToNumerics(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left: tree.Arithmetic{Left: tree.Argument{Index: 1, Type: tree.Low},
						Op:    tree.PLUS,
						Right: tree.NumericLiteral{4}},
					Rights: []tree.Numeric{tree.NumericLiteral{0}, tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t05\t1\n"+
		"ld_abs\t1C\n"+
		"add_k\t4\n"+
		"jeq_k\t01\t00\t0\n"+
		"jeq_k\t00\t01\t2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *IncludeCompilerSuite) Test_compliationOfIncludeExpressionNumericLiteralInExpressionResults(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Inclusion{
					Positive: true,
					Left:     tree.NumericLiteral{42},
					Rights: []tree.Numeric{
						tree.Arithmetic{Left: tree.Argument{Index: 1, Type: tree.Low},
							Op:    tree.PLUS,
							Right: tree.NumericLiteral{4}},
						tree.Arithmetic{Left: tree.Argument{Index: 0, Type: tree.Low},
							Op:    tree.MINUS,
							Right: tree.NumericLiteral{5}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t07\t1\n"+
		"ld_abs\t1C\n"+
		"add_k\t4\n"+
		"jeq_k\t03\t00\t2A\n"+
		"ld_abs\t14\n"+
		"sub_k\t5\n"+
		"jeq_k\t00\t01\t2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}
