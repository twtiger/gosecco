package emulator

import (
	"syscall"
	"testing"

	"github.com/twtiger/go-seccomp/data"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type EmulatorSuite struct{}

var _ = Suite(&EmulatorSuite{})

func (s *EmulatorSuite) Test_simpleReturn(c *C) {
	res := Emulate(data.SeccompData{}, []data.SockFilter{
		data.SockFilter{
			Code: syscall.BPF_RET | syscall.BPF_K,
			K:    uint32(42),
		},
	})
	c.Assert(res, Equals, uint32(42))
}
