package compiler2

// This file contains space for implementing peephole optimization
// For now it doesn't do anything, but it could once we know what
// patterns show up

// One common pattern we likely will want to fix is:
// [ST n, LDX n]
// and rewrite it into [TAX]

// We might also see [LD_IMM v, ST n, LDX n]
// This should be rewritten into [LDX_IMM v]

// We will see a lot of things like [LD_IMM v, ST n, ... LDX n, <op>]
// These, both arithmetic and comparison ones should be rewritten to use the _K
// variants. However, these is easier done in the actual compiler at the moment
// These optimizations should also check if they operate on commutative
// operators and try to put the constant to the right, where K can be used.

func (c *compilerContext) optimizeCode() {
}
