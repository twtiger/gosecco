package compiler

import (
	"testing"

	"github.com/twtiger/gosecco/asm"

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
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
					Right: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.Arithmetic{Op: tree.PLUS, Left: tree.NumericLiteral{12}, Right: tree.NumericLiteral{4}}},
				},
			},
		},
	}

	res, _ := Compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	00	0F	1\n"+
		"ld_imm	C\n"+
		"add_k	4\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	0A	0\n"+
		"ld_abs	10\n"+
		"jeq_x	07	00\n"+
		"ld_imm	C\n"+
		"add_k	4\n"+
		"tax\n"+
		"ld_abs	14\n"+
		"jeq_k	00	03	0\n"+
		"ld_abs	10\n"+
		"jeq_x	00	01\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *CompilerComparisonSuite) Test_compilationOfOrExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Or{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	04	00	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfAndExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
				},
			},
		},
	}

	res, _ := Compile(p)
	a := asm.Dump(res)

	c.Assert(a, Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	00	05	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	00	01	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfNegatedExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.Negation{
					Operand: tree.And{
						Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
						Right: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	05	00	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilationOfNestedNegatedExpression(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left: tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.Negation{
						Operand: tree.Comparison{Left: tree.Argument{Index: 1}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					},
				},
			},
		},
	}

	res, _ := Compile(p)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+ // syscallNameIndex
		"jeq_k	00	09	1\n"+ // syscall.SYS_WRITE
		"ld_abs	14\n"+ //argumentindex[0][upper]
		"jeq_k	00	07	0\n"+
		"ld_abs	10\n"+ //argumentindex[0][upper]
		"jeq_k	00	05	2A\n"+
		"ld_abs	1C\n"+ //argumentindex[1][upper]
		"jeq_k	00	03	0\n"+
		"ld_abs	18\n"+ //argumentindex[1][upper]
		"jeq_k	01	00	2A\n"+
		"ret_k	7FFF0000\n"+ //SECCOMP_RET_ALLOW
		"ret_k	0\n") //SECCOMP_RET_KILL
}

func (s *CompilerComparisonSuite) Test_compilingBooleanInsideExpressionShouldPanicSinceItsAProgrammerError(c *C) {
	p := tree.Policy{
		Rules: []tree.Rule{
			tree.Rule{
				Name: "write",
				Body: tree.And{
					Left:  tree.Comparison{Left: tree.Argument{Index: 0}, Op: tree.EQL, Right: tree.NumericLiteral{42}},
					Right: tree.BooleanLiteral{false},
				},
			},
		},
	}
	c.Assert(func() {
		Compile(p)
	}, Panics, "Programming error: there should never be any boolean literals left outside of the toplevel if the simplifier works correctly: syscall: write - (and (eq arg0 42) false)")
}
