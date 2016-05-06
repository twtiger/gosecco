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

var allow_system_call = "ret_k\t7FFF0000\n"
var kill_system_call = "ret_k\t0\n"

func (s *CompilerComparisonSuite) Test_equalsComparison_withLeftArgAndRightLiteral(c *C) {
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
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jeq_k	00	01	2A\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_comparisonWithLargest64BitNumber(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0},
					Op:    tree.EQL,
					Right: tree.NumericLiteral{9223372036854775807}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	03	7FFFFFFF\n"+
		"ld_abs	14\n"+
		"jeq_k	00	01	FFFFFFFF\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_argumentZeroNotEqualToLargest64BitNumber(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0},
					Op:    tree.NEQL,
					Right: tree.NumericLiteral{9223372036854775807}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	7FFFFFFF\n"+
		"ld_abs	14\n"+
		"jeq_k	01	00	FFFFFFFF\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_comparisonWithSecondRule(c *C) {
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
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	04	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	02	00	2A\n"+
		"ld_abs	0\n"+
		"jeq_k	00	01	99\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_greaterThanComparisonToK(c *C) {
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
		"ld_abs\t0\n"+
		"jeq_k\t00\t05\t1\n"+
		"ld_abs\t10\n"+
		"jgt_k\t00\t03\t0\n"+
		"ld_abs\t14\n"+
		"jgt_k\t00\t01\t2A\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_lessThanComparisonToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.LT, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs\t0\n"+
		"jeq_k\t00\t05\t1\n"+
		"ld_abs\t10\n"+
		"jge_k\t03\t00\t0\n"+
		"ld_abs\t14\n"+
		"jge_k\t01\t00\t2A\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_greaterThanOrEqualsToComparisonToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.GTE, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jge_k	00	01	2A\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_lessThanOrEqualsToComparisonToK(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.LTE, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	03	00	0\n"+
		"ld_abs	14\n"+
		"jge_k	01	00	2A\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_notEqualToK(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.NEQL, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	01	00	2A\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_numericLiteralGreaterThanArgument(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.GT, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jgt_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jgt_k	00	01	1\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_numericLiteralGreaterThanOrEqualsArgument(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.GTE, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jge_k	00	01	1\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_numericLiteralLessThanArgument(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.LT, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	03	00	0\n"+
		"ld_abs	14\n"+
		"jge_k	01	00	1\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_numericLiteralLessOrEqualToArgument(c *C) {
	c.Skip("p")
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.LTE, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	03	00	0\n"+
		"ld_abs	14\n"+
		"jge_k	01	00	1\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_numericLiteralNotEqualToArgument(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.NEQL, Right: tree.Argument{Index: 0}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	14\n"+
		"jeq_k	01	00	1\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_numericLiteralNotEqualToFirstArgument(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.NumericLiteral{1}, Op: tree.NEQL, Right: tree.Argument{Index: 1}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	18\n"+
		"jeq_k	00	02	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	01	00	1\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_argumentZeroNotEqualToArgumentOne(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.NEQL, Right: tree.Argument{Index: 1}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	18\n"+
		"tax\n"+
		"ld_abs	10\n"+
		"jeq_x	00	04\n"+
		"ld_abs	1C\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_x	01	00\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_argumentZeroEqualToArgumentOne(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.Argument{Index: 1}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	09	1\n"+
		"ld_abs	18\n"+
		"tax\n"+
		"ld_abs	10\n"+
		"jeq_x	00	05\n"+
		"ld_abs	1C\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_x	00	01\n"+
		allow_system_call+
		kill_system_call)
}
func (s *CompilerComparisonSuite) Test_compareArgumentToArithmeticExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 1},
					Op:    tree.NEQL,
					Right: tree.Arithmetic{Left: tree.Argument{Index: 0, Type: tree.Low}, Op: tree.PLUS, Right: tree.NumericLiteral{4}}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_abs	14\n"+
		"add_k\t4\n"+
		"tax\n"+
		"ld_abs	18\n"+
		"jeq_x	00	02\n"+
		"ld_abs	1C\n"+
		"jeq_x	01	00\n"+
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonofAtoXOfArgumentRightSideExpressionLeft(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left:  tree.Arithmetic{Left: tree.Argument{Index: 0, Type: tree.Low}, Op: tree.PLUS, Right: tree.NumericLiteral{4}},
					Op:    tree.NEQL,
					Right: tree.Argument{Index: 1}},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	08	1\n"+
		"ld_abs	14\n"+
		"add_k\t4\n"+
		"tax\n"+
		"ld_abs	18\n"+
		"jeq_x	00	02\n"+
		"ld_abs	1C\n"+
		"jeq_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonofAtoXOfExpressionsLeftAndRightSides(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{
					Left:  tree.Arithmetic{Left: tree.Argument{Index: 0, Type: tree.Low}, Op: tree.PLUS, Right: tree.NumericLiteral{10}},
					Op:    tree.NEQL,
					Right: tree.Arithmetic{Left: tree.Argument{Index: 1, Type: tree.Low}, Op: tree.MINUS, Right: tree.NumericLiteral{5}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	07	1\n"+
		"ld_abs	14\n"+
		"add_k\tA\n"+
		"tax\n"+
		"ld_abs	1C\n"+
		"sub_k\t5\n"+
		"jeq_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}
