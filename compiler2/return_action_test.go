package compiler2

import (
	"syscall"

	. "gopkg.in/check.v1"
)

type ReturnActionsSuite struct{}

var _ = Suite(&ReturnActionsSuite{})

func (s *ReturnActionsSuite) Test_returnTrap(c *C) {
	c.Assert(actionDescriptionToK("Trap"), Equals, SECCOMP_RET_TRAP)
}

func (s *ReturnActionsSuite) Test_returnKill(c *C) {
	c.Assert(actionDescriptionToK("KILL"), Equals, SECCOMP_RET_KILL)
}

func (s *ReturnActionsSuite) Test_returnTrace(c *C) {
	c.Assert(actionDescriptionToK("trace"), Equals, SECCOMP_RET_TRACE)
}

func (s *ReturnActionsSuite) Test_returnAllow(c *C) {
	c.Assert(actionDescriptionToK("AlloW"), Equals, SECCOMP_RET_ALLOW)
}

func (s *ReturnActionsSuite) Test_returnNumericValue(c *C) {
	c.Assert(actionDescriptionToK("42"), Equals, uint32(0x5002a))
}

func (s *ReturnActionsSuite) Test_returnErrName(c *C) {
	c.Assert(actionDescriptionToK("EPFNOSUPPORT"), Equals, uint32(0x50000|syscall.EPFNOSUPPORT))
}

func (s *ReturnActionsSuite) Test_returnUnknown(c *C) {
	c.Assert(actionDescriptionToK("Blarg"), Equals, uint32(0))
}
