package gosecco

import (
	"testing"

	"golang.org/x/sys/unix"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SeccompSuite struct{}

var _ = Suite(&SeccompSuite{})

func (s *SeccompSuite) Test_loadingTooBigBpf(c *C) {
	inp := make([]unix.SockFilter, 0xFFFF+1)
	res := Load(inp)
	c.Assert(res, ErrorMatches, "filter program too big: 65536 bpf instructions \\(limit = 65535\\)")
}
