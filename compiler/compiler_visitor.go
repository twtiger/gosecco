package compiler

import "github.com/twtiger/go-seccomp/tree"

const DISABLE_UNTESTED_OPTIMIZATIONS = true

type compilerVisitor struct {
	c *compiler
}

func (cv *compilerVisitor) AcceptAnd(tree.And) {}

func (cv *compilerVisitor) AcceptArgument(a tree.Argument) {
	// TODO actually load based on the index.
	// TODO sort out the messiness of the higher word etc
	cv.c.loadAt(arg0IndexLowerWord)
}

func (cv *compilerVisitor) AcceptArithmetic(tree.Arithmetic)         {}
func (cv *compilerVisitor) AcceptBinaryNegation(tree.BinaryNegation) {}
func (cv *compilerVisitor) AcceptBooleanLiteral(tree.BooleanLiteral) {}
func (cv *compilerVisitor) AcceptCall(tree.Call)                     {}

func (cv *compilerVisitor) AcceptComparison(c tree.Comparison) {
	// TODO: we can recognize whether one of the operands are
	// a literal value. This will be quite common, and in that case
	// we can actually load it directly into K and then use the K
	// version of the JMP instructions, instead of the X versions

	lit, isLit := c.Right.(tree.NumericLiteral)
	if !DISABLE_UNTESTED_OPTIMIZATIONS && isLit {
		c.Left.Accept(cv)
		switch c.Op {
		case tree.EQL:
			cv.c.jumpIfEqualTo(lit.Value, "positive", "negative")
			// TODO: deal with others here
		}
	} else {
		c.Right.Accept(cv)
		cv.c.moveAtoX()
		c.Left.Accept(cv)

		switch c.Op {
		case tree.EQL:
			cv.c.jumpIfEqualToX("positive", "negative")
			// TODO: deal with others here
		}
	}
}

func (cv *compilerVisitor) AcceptInclusion(tree.Inclusion) {}
func (cv *compilerVisitor) AcceptNegation(tree.Negation)   {}

func (cv *compilerVisitor) AcceptNumericLiteral(l tree.NumericLiteral) {
	cv.c.loadLiteral(l.Value)
}

func (cv *compilerVisitor) AcceptOr(tree.Or)             {}
func (cv *compilerVisitor) AcceptVariable(tree.Variable) {}

// func peepHole(filters []unix.SockFilter) []unix.SockFilter {
// 	one, two, three := filters[0], filters[1], filters[2]
// 	if one.Code == BPF_LD|BPF_IMM && two.Code == BPF_MISC|BPF_TAX && three.Code&(BPF_JMP|BPF_X) != 0 {
// 		return []unix.SockFilter{
// 			unix.SockFilter{
// 				Code: (three.Code & ^BPF_X) | BPF_K,
// 				Jt:   three.Jt,
// 				Jf:   three.Jf,
// 				K:    one.K,
// 			},
// 		}
// 	}
// 	return filters
// }
