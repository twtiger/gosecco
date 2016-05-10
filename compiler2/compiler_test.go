package compiler2

import (
	"syscall"

	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

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
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	01	00	1\n"+
		"ret_k	0\n"+
		"ret_k	7FFF0000\n")
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
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	0\n"+
		"jeq_k	02	00	1\n"+
		"jeq_k	01	00	99\n"+
		"ret_k	0\n"+
		"ret_k	7FFF0000\n")
}

func (s *CompilerSuite) Test_stackOverflowCreatesError(c *C) {
	cx := createCompilerContext()
	cx.stackTop = syscall.BPF_MEMWORDS
	c.Assert(cx.pushAToStack(), ErrorMatches, "the expression is too complicated to compile. Please refer to the language documentation")
}

func (s *CompilerSuite) Test_stackDoesNotOverflowRightBeforeItsLimit(c *C) {
	cx := createCompilerContext()
	cx.stackTop = syscall.BPF_MEMWORDS - 1
	c.Assert(cx.pushAToStack(), IsNil)
}

func (s *CompilerSuite) Test_stackDoesNotPopAfterReachingTheLowestIndex(c *C) {
	cx := createCompilerContext()
	cx.stackTop = 0
	c.Assert(cx.popStackToX(), ErrorMatches, "popping from empty stack - this is likely a programmer error")
}
