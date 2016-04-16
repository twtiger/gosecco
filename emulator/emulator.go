package emulator

import (
	"syscall"

	"github.com/twtiger/gosecco/data"

	"golang.org/x/sys/unix"
)

func Emulate(d data.SeccompWorkingMemory, filters []unix.SockFilter) uint32 {
	e := &emulator{data: d, filters: filters, pointer: 0}
	for {
		val, finished := e.next()
		if finished {
			return val
		}
	}
}

type emulator struct {
	data    data.SeccompWorkingMemory
	filters []unix.SockFilter
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

func (e *emulator) execRet(current unix.SockFilter) (uint32, bool) {
	switch bpfSrc(current.Code) {
	case syscall.BPF_K:
		return current.K, true
	case syscall.BPF_X:
		return e.X, true
	}
	return 0, true
}

func (e *emulator) next() (uint32, bool) {
	current := e.filters[e.pointer]
	e.pointer++

	switch bpfClass(current.Code) {
	case syscall.BPF_RET:
		return e.execRet(current)
	}

	return 0, true
}
