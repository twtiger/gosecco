package emulator

import (
	"syscall"

	"github.com/twtiger/go-seccomp/data"
)

func Emulate(d data.SeccompData, filters []data.SockFilter) uint32 {
	e := &emulator{data: d, filters: filters, pointer: 0}
	for {
		val, finished := e.next()
		if finished {
			return val
		}
	}
}

type emulator struct {
	data    data.SeccompData
	filters []data.SockFilter
	pointer uint

	X uint32
	A uint32
}

func bpfClass(code uint16) uint16 {
	return code & 0x07
}

func bpfSize(code uint16) uint16 {
	return code & 0x18
}

func bpfMode(code uint16) uint16 {
	return code & 0xe0
}

func bpfOp(code uint16) uint16 {
	return code & 0xf0
}

func bpfSrc(code uint16) uint16 {
	return code & 0x08
}

func (e *emulator) next() (uint32, bool) {
	current := e.filters[e.pointer]
	e.pointer++

	switch bpfClass(current.Code) {
	case syscall.BPF_RET:
		switch bpfSrc(current.Code) {
		case syscall.BPF_K:
			return current.K, true
		case syscall.BPF_X:
			return e.X, true
		}
	}

	return 0, true
}
