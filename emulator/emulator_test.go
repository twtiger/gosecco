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

func (s *EmulatorSuite) Test_simpleReturnK(c *C) {
	res := Emulate(data.SeccompData{}, []data.SockFilter{
		data.SockFilter{
			Code: syscall.BPF_RET | syscall.BPF_K,
			K:    uint32(42),
		},
	})
	c.Assert(res, Equals, uint32(42))
}

func (s *EmulatorSuite) Test_simpleReturnX(c *C) {
	e := &emulator{
		data: data.SeccompData{},
		filters: []data.SockFilter{
			data.SockFilter{
				Code: syscall.BPF_RET | syscall.BPF_X,
				K:    uint32(42),
			},
		},
		pointer: 0,

		X: uint32(23),
	}

	res, _ := e.next()

	c.Assert(res, Equals, uint32(23))
}
