package compiler2

import (
	"github.com/twtiger/gosecco/asm"
	. "gopkg.in/check.v1"
)

type PrefixSuite struct{}

var _ = Suite(&PrefixSuite{})

func (s *PrefixSuite) Test_compilesAuditArch(c *C) {
	c.Skip("waiting on actions")
	ctx := createCompilerContext()
	ctx.compileAuditArchCheck(label("TODO FIX ME"))
	c.Assert(asm.Dump(ctx.result), Equals, "")

}
