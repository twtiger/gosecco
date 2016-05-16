package compiler

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// This file contains space for implementing peephole optimization
// For now it doesn't do anything, but it could once we know what
// patterns show up

// One common pattern we likely will want to fix is:
// [ST n, LDX n]
// and rewrite it into [TAX]

// We might also see [LD_IMM v, ST n, LDX n]
// This should be rewritten into [LDX_IMM v]

// Some patterns look amenable to optimization but in practice won't be
// - it's important that we are wary of trying to fix up jumps too much.

// We will see a lot of things like [LD_IMM v, ST n, ... LDX n, <op>]
// These, both arithmetic and comparison ones should be rewritten to use the _K
// variants. However, these is easier done in the actual compiler at the moment
// These optimizations should also check if they operate on commutative
// operators and try to put the constant to the right, where K can be used.

func (c *compilerContext) optimizeCode() {
	index := 0

	// Do not pull out the length calculation here, since the length
	// can change during optimization
	for index < len(c.result) {
		c.optimizeAt(index)
		index++
	}
}

type optimizer func(*compilerContext, int) bool

var optimizers = []optimizer{
	jumpAfterConditionalJumpOptimizer,
}

func (c *compilerContext) optimizeAt(i int) {
	for _, o := range optimizers {
		o(c, i)
	}
}

func isJump(s unix.SockFilter) bool {
	return bpfClass(s.Code) == syscall.BPF_JMP
}

func isConditionalJump(s unix.SockFilter) bool {
	return isJump(s) && bpfOp(s.Code) != syscall.BPF_JA
}

func isUnconditionalJump(s unix.SockFilter) bool {
	return isJump(s) && bpfOp(s.Code) == syscall.BPF_JA
}

// hasJumpTarget will return true if the conditional jump given has at least one of
// its target being the potential target.
func hasJumpTarget(c *compilerContext, conditionalIndex, potentialTarget int) bool {
	return c.jts.hasJumpTarget(c, conditionalIndex, potentialTarget) ||
		c.jfs.hasJumpTarget(c, conditionalIndex, potentialTarget)
}

// hasOnlyJumpFrom will make sure that the given jump target is only jumped to from the
// given expected from location - this will return false if the expectedFrom is a
// conditional jump where both conditions point to the jump target
func hasOnlyJumpFrom(c *compilerContext, jumpTarget, expectedFrom int) bool {
	return (c.jts.countJumpsFrom(c, jumpTarget, expectedFrom)+
		c.jfs.countJumpsFrom(c, jumpTarget, expectedFrom)+
		c.uconds.countJumpsFrom(c, jumpTarget, expectedFrom)) == 1 &&
		(c.jts.countJumpsFromAny(c, jumpTarget)+
			c.jfs.countJumpsFromAny(c, jumpTarget)+
			c.uconds.countJumpsFromAny(c, jumpTarget)) == 1
}

// isNotOversizedJump takes a pointer to an unconditional jump and return true if
// the jump is smaller than the max jump size.
func isNotOversizedJump(c *compilerContext, jumpPoint int) bool {
	return !c.uconds.jumpSizeIsOversized(c, jumpPoint)
}

func redirectJumpOf(c *compilerContext, ucond, cond int) {
	newJumpTarget := c.uconds.jumpTargetOf(ucond)
	c.uconds.removeJumpTarget(ucond)
	sourceJm := c.jts
	if c.jfs.countJumpsFrom(c, cond, ucond) == 1 {
		sourceJm = c.jfs
	}
	oldLabel := c.labels.labelsAt(ucond)[0]
	sourceJm.redirectJump(oldLabel, newJumpTarget)
	c.labels.removeLabel(oldLabel)
}

func (c *compilerContext) removeInstructionAt(index int) {
	c.result = append(c.result[:index], c.result[index+1:]...)
}

// jumpAfterConditionalJumpOptimizer will optimize situations where a JMP instruction
// directly follows a conditional jump where one of the arms of the conditional jump
// is zero. It will make sure that no other jump points end up on the specific JMP instruction
// before removing it. It will also make sure the resulting jump is not too large.
// An example of a fragment that would be changed would be this:
//    jeq_k	00	01	3D
//    jmp	13
// This can be optimized to:
//    jeq_k	13	00	3D
func jumpAfterConditionalJumpOptimizer(c *compilerContext, ix int) bool {
	// Found a good jump at: 3 for:
	// ld_abs	4
	// jeq_k	00	00	C000003E
	// ld_abs	0
	// jeq_k	00	00	1
	// jmp	0
	// jeq_k	00	00	0
	// jmp	0
	// jeq_k	00	00	10
	// jmp	0
	// jeq_k	00	00	13E
	// jmp	0
	// jmp	0
	// ret_k	7FFF0000
	// ret_k	0

	// jts := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel001": []int{1},
	// 	"generatedLabel004": []int{3},
	// 	"generatedLabel006": []int{5},
	// 	"generatedLabel008": []int{7},
	// 	"generatedLabel010": []int{9},
	// },
	// 	positionToLabel: map[int]compiler.label{
	// 		1: "generatedLabel001",
	// 		3: "generatedLabel004",
	// 		5: "generatedLabel006",
	// 		7: "generatedLabel008",
	// 		9: "generatedLabel010",
	// 	}}

	// jfs := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel000": []int{1},
	// 	"generatedLabel002": []int{3},
	// 	"generatedLabel005": []int{5},
	// 	"generatedLabel007": []int{7},
	// 	"generatedLabel009": []int{9},
	// },
	// 	positionToLabel: map[int]compiler.label{
	// 		1: "generatedLabel000",
	// 		3: "generatedLabel002",
	// 		5: "generatedLabel005",
	// 		7: "generatedLabel007",
	// 		9: "generatedLabel009",
	// 	}}

	// uco := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel000": []int{11},
	// 	"generatedLabel003": []int{4, 6, 8, 10},
	// },
	// 	positionToLabel: map[int]compiler.label{
	// 		10: "generatedLabel003",
	// 		11: "generatedLabel000",
	// 		4:  "generatedLabel003",
	// 		6:  "generatedLabel003",
	// 		8:  "generatedLabel003",
	// 	}}

	// labels := &compiler.labelMap{labelToPosition: map[compiler.label]int{
	// 	"generatedLabel000": 13,
	// 	"generatedLabel001": 2,
	// 	"generatedLabel002": 5,
	// 	"generatedLabel003": 12,
	// 	"generatedLabel004": 4,
	// 	"generatedLabel005": 7,
	// 	"generatedLabel006": 6,
	// 	"generatedLabel007": 9,
	// 	"generatedLabel008": 8,
	// 	"generatedLabel009": 11,
	// 	"generatedLabel010": 10,
	// }, positionToLabel: map[int][]compiler.label{
	// 	10: []compiler.label{"generatedLabel010"},
	// 	11: []compiler.label{"generatedLabel009"},
	// 	12: []compiler.label{"generatedLabel003"},
	// 	13: []compiler.label{"generatedLabel000"},
	// 	2:  []compiler.label{"generatedLabel001"},
	// 	4:  []compiler.label{"generatedLabel004"},
	// 	5:  []compiler.label{"generatedLabel002"},
	// 	6:  []compiler.label{"generatedLabel006"},
	// 	7:  []compiler.label{"generatedLabel005"},
	// 	8:  []compiler.label{"generatedLabel008"},
	// 	9:  []compiler.label{"generatedLabel007"},
	// }}

	// //    after:
	// // ld_abs	4
	// // jeq_k	00	00	C000003E
	// // ld_abs	0
	// // jeq_k	00	00	1
	// // jeq_k	00	00	0
	// // jmp	0
	// // jeq_k	00	00	10
	// // jmp	0
	// // jeq_k	00	00	13E
	// // jmp	0
	// // jmp	0
	// // ret_k	7FFF0000
	// // ret_k	0

	// jts := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel001": []int{1},
	// 	"generatedLabel003": []int{3},
	// 	"generatedLabel006": []int{4},
	// 	"generatedLabel008": []int{6},
	// 	"generatedLabel010": []int{8},
	// }, positionToLabel: map[int]compiler.label{
	// 	1: "generatedLabel001",
	// 	3: "generatedLabel003",
	// 	4: "generatedLabel006",
	// 	6: "generatedLabel008",
	// 	8: "generatedLabel010",
	// }}

	// jfs := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel000": []int{1},
	// 	"generatedLabel002": []int{3},
	// 	"generatedLabel005": []int{4},
	// 	"generatedLabel007": []int{6},
	// 	"generatedLabel009": []int{8},
	// }, positionToLabel: map[int]compiler.label{
	// 	1: "generatedLabel000",
	// 	3: "generatedLabel002",
	// 	4: "generatedLabel005",
	// 	6: "generatedLabel007",
	// 	8: "generatedLabel009",
	// }}

	// uco := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel000": []int{10},
	// 	"generatedLabel003": []int{4, 5, 7, 9},
	// },
	// 	positionToLabel: map[int]compiler.label{
	// 		10: "generatedLabel000",
	// 		4:  "generatedLabel003",
	// 		5:  "generatedLabel003",
	// 		7:  "generatedLabel003",
	// 		9:  "generatedLabel003",
	// 	}}

	// labels := &compiler.labelMap{labelToPosition: map[compiler.label]int{
	// 	"generatedLabel000": 12,
	// 	"generatedLabel001": 2,
	// 	"generatedLabel002": 4,
	// 	"generatedLabel003": 11,
	// 	"generatedLabel005": 6,
	// 	"generatedLabel006": 5,
	// 	"generatedLabel007": 8,
	// 	"generatedLabel008": 7,
	// 	"generatedLabel009": 10,
	// 	"generatedLabel010": 9,
	// }, positionToLabel: map[int][]compiler.label{2: []compiler.label{"generatedLabel001"}, 12: []compiler.label{"generatedLabel000"}, 11: []compiler.label{"generatedLabel003"}, 10: []compiler.label{"generatedLabel009"}, 6: []compiler.label{"generatedLabel005"}, 9: []compiler.label{"generatedLabel010"}, 4: []compiler.label{"generatedLabel002"}, 5: []compiler.label{"generatedLabel006"}, 7: []compiler.label{"generatedLabel008"}, 8: []compiler.label{"generatedLabel007"}}}

	// jts := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel001": []int{1},
	// 	"generatedLabel003": []int{6},
	// }, positionToLabel: map[int]compiler.label{1: "generatedLabel001", 6: "generatedLabel003"}}

	// jfs := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel002": []int{3},
	// 	"generatedLabel005": []int{4},
	// 	"generatedLabel007": []int{5},
	// 	"generatedLabel009": []int{6},
	// 	"generatedLabel000": []int{1},
	// }, positionToLabel: map[int]compiler.label{4: "generatedLabel005", 5: "generatedLabel007", 6: "generatedLabel009", 1: "generatedLabel000", 3: "generatedLabel002"}}

	// uco := &compiler.jumpMap{labelToPosition: map[compiler.label][]int{
	// 	"generatedLabel000": []int{7},
	// }, positionToLabel: map[int]compiler.label{7: "generatedLabel000"}}

	// labels := &compiler.labelMap{labelToPosition: map[compiler.label]int{
	// 	"generatedLabel000": 9,
	// 	"generatedLabel001": 2,
	// 	"generatedLabel002": 4,
	// 	"generatedLabel003": 8,
	// 	"generatedLabel005": 5,
	// 	"generatedLabel007": 6,
	// 	"generatedLabel009": 7,
	// }, positionToLabel: map[int][]compiler.label{9: []compiler.label{"generatedLabel000"}, 5: []compiler.label{"generatedLabel005"}, 2: []compiler.label{"generatedLabel001"}, 4: []compiler.label{"generatedLabel002"}, 6: []compiler.label{"generatedLabel007"}, 8: []compiler.label{"generatedLabel003"}, 7: []compiler.label{"generatedLabel009"}}}

	if ix+1 < len(c.result) {
		oneIndex, twoIndex := ix, ix+1
		one, two := c.result[oneIndex], c.result[twoIndex]
		if isConditionalJump(one) &&
			isUnconditionalJump(two) &&
			hasJumpTarget(c, oneIndex, twoIndex) &&
			hasOnlyJumpFrom(c, twoIndex, oneIndex) &&
			isNotOversizedJump(c, twoIndex) {

			// fmt.Printf("Found a good jump at: %d for: \n%s\n\njts: %#v\n\njfs: %#v\n\nuco: %#v\n\nlabels: %#v\n\n", ix, asm.Dump(c.result), c.jts, c.jfs, c.uconds, c.labels)

			redirectJumpOf(c, twoIndex, oneIndex)
			c.shiftJumpsBy(oneIndex+1, -1)
			c.removeInstructionAt(twoIndex)

			// fmt.Printf("   after: \n%s\n\njts: %#v\n\njfs: %#v\n\nuco: %#v\n\nlabels: %#v\n\n", asm.Dump(c.result), c.jts, c.jfs, c.uconds, c.labels)
		}
	}

	return false
}
