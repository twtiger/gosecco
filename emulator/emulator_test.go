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

func (s *EmulatorSuite) Test_loadValues(c *C) {
	e := &emulator{
		data: data.SeccompWorkingMemory{Arch: 15, Args: [6]uint64{0, 0, 0, 12423423, 0, 0}},
		filters: []unix.SockFilter{
			unix.SockFilter{
				Code: syscall.BPF_LD | syscall.BPF_W | syscall.BPF_ABS,
				K:    uint32(4),
			},
			unix.SockFilter{
				Code: syscall.BPF_LD | syscall.BPF_W | syscall.BPF_ABS,
				K:    uint32(44),
			},
			unix.SockFilter{
				Code: syscall.BPF_LD | syscall.BPF_W | syscall.BPF_IND,
				K:    uint32(2),
			},
			unix.SockFilter{
				Code: syscall.BPF_LD | syscall.BPF_W | syscall.BPF_IND,
				K:    uint32(42),
			},
			unix.SockFilter{
				Code: syscall.BPF_LD | syscall.BPF_W | syscall.BPF_LEN,
			},
			unix.SockFilter{
				Code: syscall.BPF_LD | syscall.BPF_IMM,
				K:    uint32(23),
			},
		},
		pointer: 0,
		X:       uint32(2),
	}

	e.next()

	c.Assert(e.A, Equals, uint32(15))

	e.next()

	c.Assert(e.A, Equals, uint32(12423423))

	e.next()

	c.Assert(e.A, Equals, uint32(15))

	e.next()

	c.Assert(e.A, Equals, uint32(12423423))

	e.next()

	c.Assert(e.A, Equals, uint32(64))

	e.next()

	c.Assert(e.A, Equals, uint32(23))
}

func (s *EmulatorSuite) Test_loadValuesIntoX(c *C) {
	e := &emulator{
		data: data.SeccompWorkingMemory{},
		filters: []unix.SockFilter{
			unix.SockFilter{
				Code: syscall.BPF_LDX | syscall.BPF_IMM,
				K:    uint32(234),
			},
			unix.SockFilter{
				Code: syscall.BPF_LDX | syscall.BPF_W | syscall.BPF_LEN,
			},
		},
		pointer: 0,
	}

	e.next()

	c.Assert(e.X, Equals, uint32(234))

	e.next()

	c.Assert(e.X, Equals, uint32(64))
}

func aluAndK(c *C, op uint16, a, k, expected uint32) {
	e := &emulator{
		data: data.SeccompWorkingMemory{},
		filters: []unix.SockFilter{
			unix.SockFilter{
				Code: syscall.BPF_ALU | syscall.BPF_K | op,
				K:    k,
			},
		},
		pointer: 0,
		A:       a,
	}

	e.next()

	c.Assert(e.A, Equals, expected)
}

func aluAndX(c *C, op uint16, a, k, expected uint32) {
	e := &emulator{
		data: data.SeccompWorkingMemory{},
		filters: []unix.SockFilter{
			unix.SockFilter{
				Code: syscall.BPF_ALU | syscall.BPF_X | op,
			},
		},
		pointer: 0,
		A:       a,
		X:       k,
	}

	e.next()

	c.Assert(e.A, Equals, expected)
}

func (s *EmulatorSuite) Test_aluAandK(c *C) {
	aluAndK(c, syscall.BPF_ADD, 15, 42, 57)
	aluAndK(c, syscall.BPF_SUB, 10, 3, 7)
	aluAndK(c, syscall.BPF_MUL, 10, 3, 30)
	aluAndK(c, syscall.BPF_DIV, 10, 3, 3)
	aluAndK(c, syscall.BPF_AND, 32425, 1211, 32425&1211)
	aluAndK(c, syscall.BPF_OR, 32425, 1211, 32425|1211)
	aluAndK(c, BPF_XOR, 32425, 1211, 32425^1211)
	aluAndK(c, syscall.BPF_LSH, 10, 3, 80)
	aluAndK(c, syscall.BPF_RSH, 80, 3, 10)
	aluAndK(c, BPF_MOD, 10, 3, 1)
	aluAndK(c, syscall.BPF_NEG, 80, 0, 0xFFFFFFB0)
}

func (s *EmulatorSuite) Test_aluAandX(c *C) {
	aluAndX(c, syscall.BPF_ADD, 15, 42, 57)
	aluAndX(c, syscall.BPF_SUB, 10, 3, 7)
	aluAndX(c, syscall.BPF_MUL, 10, 3, 30)
	aluAndX(c, syscall.BPF_DIV, 10, 3, 3)
	aluAndX(c, syscall.BPF_AND, 32425, 1211, 32425&1211)
	aluAndX(c, syscall.BPF_OR, 32425, 1211, 32425|1211)
	aluAndX(c, BPF_XOR, 32425, 1211, 32425^1211)
	aluAndX(c, syscall.BPF_LSH, 10, 3, 80)
	aluAndX(c, syscall.BPF_RSH, 80, 3, 10)
	aluAndX(c, BPF_MOD, 10, 3, 1)
}

func (s *EmulatorSuite) Test_misc(c *C) {
	e := &emulator{
		data: data.SeccompWorkingMemory{},
		filters: []unix.SockFilter{
			unix.SockFilter{
				Code: syscall.BPF_MISC | syscall.BPF_TAX,
			},
		},
		pointer: 0,
		A:       42,
		X:       23,
	}

	e.next()

	c.Assert(e.X, Equals, uint32(42))

	e = &emulator{
		data: data.SeccompWorkingMemory{},
		filters: []unix.SockFilter{
			unix.SockFilter{
				Code: syscall.BPF_MISC | syscall.BPF_TXA,
			},
		},
		pointer: 0,
		A:       42,
		X:       23,
	}

	e.next()

	c.Assert(e.A, Equals, uint32(23))
}
