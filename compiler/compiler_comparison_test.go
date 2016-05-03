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

func (s *CompilerComparisonSuite) Test_compilationOfEqualsComparison_withLeftArgAndRightLit(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
			},
		},
	}

	res, _ := Compile(p)

	allow_system_call := "ret_k\t7FFF0000\n"
	kill_system_call := "ret_k\t0\n"

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+ // is value in accumulator the system call write, where 1 is write?
		"ld_abs	10\n"+ // load upper half of arg0 into A
		"jeq_k	00	03	0\n"+ // compare what is in A to upper half of our constant, 42 (which is also 64 bits)
		"ld_abs	14\n"+ // load the lower half of argument 0 which is a 32 bit value into A
		"jeq_k	00	01	2A\n"+ // compare what is in A to lower half of our constant, 42 (which is also 64 bits)
		allow_system_call+
		kill_system_call)
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonWithLargerNumber(c *C) {
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
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
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
		"ld_abs	10\n"+
		"jgt_k	00	03	0\n"+
		"ld_abs	14\n"+
		"jgt_k	00	01	2A\n"+
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
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jgt_k	03	00	0\n"+
		"ld_abs	14\n"+
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
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	00	03	0\n"+
		"ld_abs	14\n"+
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
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jge_k	03	00	0\n"+
		"ld_abs	14\n"+
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
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	10\n"+
		"jeq_k	03	00	0\n"+
		"ld_abs	14\n"+
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterThanRightSide(c *C) {
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
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfGreaterOrEqualsToRightSide(c *C) {
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
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfLessThanKLeftSide(c *C) {
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
		"jgt_k	03	00	0\n"+
		"ld_abs	14\n"+
		"jgt_k	01	00	1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfLessOrEqualsToKLeftSide(c *C) {
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
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfNotEqualsKLeftSide(c *C) {
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
		"jeq_k	03	00	0\n"+
		"ld_abs	14\n"+
		"jeq_k	01	00	1\n"+
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
	a := asm.Dump(res)
	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	05	1\n"+
		"ld_abs	18\n"+
		"jeq_k	03	00	0\n"+
		"ld_abs	1C\n"+
		"jeq_k	01	00	1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonofAtoXOfTwoArguments(c *C) {
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
		"jeq_x	05	00\n"+
		"ld_abs	1C\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfComparisonofAtoXOfArgumentLeftSideExpressionRight(c *C) {
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
		"jeq_x	03	00\n"+
		"ld_abs	1C\n"+
		"jeq_x	01	00\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
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
		"jeq_x	03	00\n"+
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
