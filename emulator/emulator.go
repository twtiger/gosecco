package emulator

import (
	"fmt"
	"log"
	"syscall"

	"github.com/twtiger/gosecco/data"

	"golang.org/x/sys/unix"
)

func init() {
	log.SetFlags(0)
}

// Emulate will execute a seccomp filter program against the given working memory.
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
	pointer uint32

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
func bpfMiscOp(code uint16) uint16 {
	return code & 0xf8
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
	default:
		panic(fmt.Sprintf("Invalid ret source: %d", bpfSrc(current.Code)))
	}
	return 0, true
}

func (e *emulator) getFromWorkingMemory(ix uint32) uint32 {
	switch ix {
	case 0:
		return uint32(e.data.NR)
	case 4:
		return e.data.Arch
	case 8:
		return uint32(e.data.InstructionPointer >> 32)
	case 12:
		return uint32(e.data.InstructionPointer & 0xFFFF)
	case 16:
		return uint32(e.data.Args[0] >> 32)
	case 20:
		return uint32(e.data.Args[0] & 0xFFFFFFFF)
	case 24:
		return uint32(e.data.Args[1] >> 32)
	case 28:
		return uint32(e.data.Args[1] & 0xFFFFFFFF)
	case 32:
		return uint32(e.data.Args[2] >> 32)
	case 36:
		return uint32(e.data.Args[2] & 0xFFFFFFFF)
	case 40:
		return uint32(e.data.Args[3] >> 32)
	case 44:
		return uint32(e.data.Args[3] & 0xFFFFFFFF)
	case 48:
		return uint32(e.data.Args[4] >> 32)
	case 52:
		return uint32(e.data.Args[4] & 0xFFFFFFFF)
	case 56:
		return uint32(e.data.Args[5] >> 32)
	case 60:
		return uint32(e.data.Args[5] & 0xFFFFFFFF)
	default:
		return 0
	}
}

func (e *emulator) loadFromWorkingMemory(ix uint32) {
	e.A = e.getFromWorkingMemory(ix)
}

func (e *emulator) execLd(current unix.SockFilter) (uint32, bool) {
	cd := current.Code

	if bpfSize(cd) != syscall.BPF_W {
		panic("Invalid code, we can't load smaller values than wide ones")
	}

	switch bpfMode(cd) {
	case syscall.BPF_ABS:
		e.loadFromWorkingMemory(current.K)
	case syscall.BPF_IND:
		e.loadFromWorkingMemory(e.X + current.K)
	case syscall.BPF_LEN:
		e.A = uint32(64)
	case syscall.BPF_IMM:
		e.A = current.K
	default:
		panic(fmt.Sprintf("Invalid mode: %d", bpfMode(cd)))
	}
	return 0, false
}

func (e *emulator) execLdx(current unix.SockFilter) (uint32, bool) {
	cd := current.Code

	if bpfSize(cd) != syscall.BPF_W {
		panic("Invalid code, we can't load smaller values than wide ones")
	}

	switch bpfMode(cd) {
	case syscall.BPF_LEN:
		e.X = uint32(64)
	case syscall.BPF_IMM:
		e.X = current.K
	default:
		panic(fmt.Sprintf("Invalid mode: %d", bpfMode(cd)))
	}
	return 0, false
}

// BPF_MOD is BPF_MOD - it is supported in Linux from v3.7+, but not in go's syscall...
const BPF_MOD = 0x90

// BPF_XOR is BPF_XOR - it is supported in Linux from v3.7+, but not in go's syscall...
const BPF_XOR = 0xa0

func (e *emulator) execAlu(current unix.SockFilter) (uint32, bool) {
	cd := current.Code

	right := uint32(0)

	switch bpfSrc(cd) {
	case syscall.BPF_K:
		right = current.K
	case syscall.BPF_X:
		right = e.X
	default:
		panic(fmt.Sprintf("Invalid source for right hand side of operation: %d", bpfSrc(cd)))
	}

	switch bpfOp(cd) {
	case syscall.BPF_ADD:
		e.A += right
	case syscall.BPF_SUB:
		e.A -= right
	case syscall.BPF_MUL:
		e.A *= right
	case syscall.BPF_DIV:
		e.A /= right
	case syscall.BPF_AND:
		e.A &= right
	case syscall.BPF_OR:
		e.A |= right
	case BPF_XOR:
		e.A ^= right
	case syscall.BPF_LSH:
		e.A <<= right
	case syscall.BPF_RSH:
		e.A >>= right
	case BPF_MOD:
		e.A %= right
	case syscall.BPF_NEG:
		e.A = -e.A
	default:
		panic(fmt.Sprintf("Invalid op: %d", bpfOp(cd)))
	}
	return 0, false
}

func (e *emulator) execMisc(current unix.SockFilter) (uint32, bool) {
	cd := current.Code

	switch bpfMiscOp(cd) {
	case syscall.BPF_TAX:
		e.X = e.A
	case syscall.BPF_TXA:
		e.A = e.X
	default:
		panic(fmt.Sprintf("Invalid op: %d", bpfMiscOp(cd)))
	}
	return 0, false
}

func (e *emulator) execJmp(current unix.SockFilter) (uint32, bool) {
	cd := current.Code

	right := uint32(0)
	switch bpfSrc(cd) {
	case syscall.BPF_K:
		right = current.K
	case syscall.BPF_X:
		right = e.X
	default:
		panic(fmt.Sprintf("Invalid source for right hand side of operation: %d", bpfSrc(cd)))
	}

	switch bpfOp(cd) {
	case syscall.BPF_JA:
		e.pointer += current.K
	case syscall.BPF_JGT:
		if e.A > right {
			e.pointer += uint32(current.Jt)
		} else {
			e.pointer += uint32(current.Jf)
		}
	case syscall.BPF_JGE:
		if e.A >= right {
			e.pointer += uint32(current.Jt)
		} else {
			e.pointer += uint32(current.Jf)
		}
	case syscall.BPF_JEQ:
		if e.A == right {
			e.pointer += uint32(current.Jt)
		} else {
			e.pointer += uint32(current.Jf)
		}
	case syscall.BPF_JSET:
		if e.A&right != 0 {
			e.pointer += uint32(current.Jt)
		} else {
			e.pointer += uint32(current.Jf)
		}
	default:
		panic(fmt.Sprintf("Invalid op: %d", bpfOp(cd)))
	}
	return 0, false
}

func (e *emulator) next() (uint32, bool) {
	if e.pointer >= uint32(len(e.filters)) {
		return 0, true
	}

	current := e.filters[e.pointer]
	//	log.Printf("%03d:  %s\n", e.pointer, formatInstruction(current))
	e.pointer++
	switch bpfClass(current.Code) {
	case syscall.BPF_RET:
		return e.execRet(current)
	case syscall.BPF_LD:
		return e.execLd(current)
	case syscall.BPF_LDX:
		return e.execLdx(current)
	case syscall.BPF_ALU:
		return e.execAlu(current)
	case syscall.BPF_MISC:
		return e.execMisc(current)
	case syscall.BPF_JMP:
		return e.execJmp(current)
	}

	return 0, true
}

func instructionName(code uint16) string {
	switch code {
	case syscall.BPF_RET | syscall.BPF_K:
		return "ret_k"
	case syscall.BPF_RET | syscall.BPF_X:
		return "ret_x"
	case syscall.BPF_LD | syscall.BPF_W | syscall.BPF_ABS:
		return "ld_abs"
	case syscall.BPF_LD | syscall.BPF_W | syscall.BPF_IND:
		return "ld_ind"
	case syscall.BPF_LD | syscall.BPF_W | syscall.BPF_LEN:
		return "ld_len"
	case syscall.BPF_LD | syscall.BPF_W | syscall.BPF_IMM:
		return "ld_imm"
	case syscall.BPF_LDX | syscall.BPF_W | syscall.BPF_LEN:
		return "ldx_len"
	case syscall.BPF_LDX | syscall.BPF_W | syscall.BPF_IMM:
		return "ldx_imm"
	case syscall.BPF_ALU | syscall.BPF_ADD | syscall.BPF_K:
		return "add_k"
	case syscall.BPF_ALU | syscall.BPF_ADD | syscall.BPF_X:
		return "add_x"
	case syscall.BPF_ALU | syscall.BPF_SUB | syscall.BPF_K:
		return "sub_k"
	case syscall.BPF_ALU | syscall.BPF_SUB | syscall.BPF_X:
		return "sub_x"
	case syscall.BPF_ALU | syscall.BPF_MUL | syscall.BPF_K:
		return "mul_k"
	case syscall.BPF_ALU | syscall.BPF_MUL | syscall.BPF_X:
		return "mul_x"
	case syscall.BPF_ALU | syscall.BPF_DIV | syscall.BPF_K:
		return "div_k"
	case syscall.BPF_ALU | syscall.BPF_DIV | syscall.BPF_X:
		return "div_x"
	case syscall.BPF_ALU | syscall.BPF_AND | syscall.BPF_K:
		return "and_k"
	case syscall.BPF_ALU | syscall.BPF_AND | syscall.BPF_X:
		return "and_x"
	case syscall.BPF_ALU | syscall.BPF_OR | syscall.BPF_K:
		return "or_k"
	case syscall.BPF_ALU | syscall.BPF_OR | syscall.BPF_X:
		return "or_x"
	case syscall.BPF_ALU | BPF_XOR | syscall.BPF_K:
		return "xor_k"
	case syscall.BPF_ALU | BPF_XOR | syscall.BPF_X:
		return "xor_x"
	case syscall.BPF_ALU | syscall.BPF_LSH | syscall.BPF_K:
		return "lsh_k"
	case syscall.BPF_ALU | syscall.BPF_LSH | syscall.BPF_X:
		return "lsh_x"
	case syscall.BPF_ALU | syscall.BPF_RSH | syscall.BPF_K:
		return "rsh_k"
	case syscall.BPF_ALU | syscall.BPF_RSH | syscall.BPF_X:
		return "rsh_x"
	case syscall.BPF_ALU | BPF_MOD | syscall.BPF_K:
		return "mod_k"
	case syscall.BPF_ALU | BPF_MOD | syscall.BPF_X:
		return "mod_x"
	case syscall.BPF_ALU | syscall.BPF_NEG:
		return "neg"
	case syscall.BPF_MISC | syscall.BPF_TAX:
		return "tax"
	case syscall.BPF_MISC | syscall.BPF_TXA:
		return "txa"
	case syscall.BPF_JMP | syscall.BPF_JA:
		return "jmp"
	case syscall.BPF_JMP | syscall.BPF_JGT | syscall.BPF_K:
		return "jgt_k"
	case syscall.BPF_JMP | syscall.BPF_JGE | syscall.BPF_K:
		return "jge_k"
	case syscall.BPF_JMP | syscall.BPF_JEQ | syscall.BPF_K:
		return "jeq_k"
	case syscall.BPF_JMP | syscall.BPF_JSET | syscall.BPF_K:
		return "jset_k"
	case syscall.BPF_JMP | syscall.BPF_JGT | syscall.BPF_X:
		return "jgt_x"
	case syscall.BPF_JMP | syscall.BPF_JGE | syscall.BPF_X:
		return "jge_x"
	case syscall.BPF_JMP | syscall.BPF_JEQ | syscall.BPF_X:
		return "jeq_x"
	case syscall.BPF_JMP | syscall.BPF_JSET | syscall.BPF_X:
		return "jset_x"
	default:
		return ""
	}
}

func formatInstruction(current unix.SockFilter) string {
	return fmt.Sprintf("%04x  %-10s   +%02x    -%02x    %08x", current.Code, instructionName(current.Code), current.Jt, current.Jf, current.K)
}
