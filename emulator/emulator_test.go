package emulator

import (
	"syscall"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/twtiger/gosecco/data"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type EmulatorSuite struct{}

var _ = Suite(&EmulatorSuite{})

func (s *EmulatorSuite) Test_simpleReturnK(c *C) {
	res := Emulate(data.SeccompWorkingMemory{}, []unix.SockFilter{
		unix.SockFilter{
			Code: syscall.BPF_RET | syscall.BPF_K,
			K:    uint32(42),
		},
	})
	c.Assert(res, Equals, uint32(42))
}

func (s *EmulatorSuite) Test_simpleReturnX(c *C) {
	e := &emulator{
		data: data.SeccompWorkingMemory{},
		filters: []unix.SockFilter{
			unix.SockFilter{
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
