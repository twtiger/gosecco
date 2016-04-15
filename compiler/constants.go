package compiler

import "syscall"

const (
	BPF_A             = syscall.BPF_A
	BPF_ABS           = syscall.BPF_ABS
	BPF_ADD           = syscall.BPF_ADD
	BPF_ALU           = syscall.BPF_ALU
	BPF_AND           = syscall.BPF_AND
	BPF_B             = syscall.BPF_B
	BPF_DIV           = syscall.BPF_DIV
	BPF_H             = syscall.BPF_H
	BPF_IMM           = syscall.BPF_IMM
	BPF_IND           = syscall.BPF_IND
	BPF_JA            = syscall.BPF_JA
	BPF_JEQ           = syscall.BPF_JEQ
	BPF_JGE           = syscall.BPF_JGE
	BPF_JGT           = syscall.BPF_JGT
	BPF_JMP           = syscall.BPF_JMP
	BPF_JSET          = syscall.BPF_JSET
	BPF_K             = syscall.BPF_K
	BPF_LD            = syscall.BPF_LD
	BPF_LDX           = syscall.BPF_LDX
	BPF_LEN           = syscall.BPF_LEN
	BPF_LSH           = syscall.BPF_LSH
	BPF_MAJOR_VERSION = syscall.BPF_MAJOR_VERSION
	BPF_MAXINSNS      = syscall.BPF_MAXINSNS
	BPF_MEM           = syscall.BPF_MEM
	BPF_MEMWORDS      = syscall.BPF_MEMWORDS
	BPF_MINOR_VERSION = syscall.BPF_MINOR_VERSION
	BPF_MISC          = syscall.BPF_MISC
	BPF_MSH           = syscall.BPF_MSH
	BPF_MUL           = syscall.BPF_MUL
	BPF_NEG           = syscall.BPF_NEG
	BPF_OR            = syscall.BPF_OR
	BPF_RET           = syscall.BPF_RET
	BPF_RSH           = syscall.BPF_RSH
	BPF_ST            = syscall.BPF_ST
	BPF_STX           = syscall.BPF_STX
	BPF_SUB           = syscall.BPF_SUB
	BPF_TAX           = syscall.BPF_TAX
	BPF_TXA           = syscall.BPF_TXA
	BPF_W             = syscall.BPF_W
	BPF_X             = syscall.BPF_X

	SECCOMP_RET_KILL  = uint32(0x00000000) /* kill the task immediately */
	SECCOMP_RET_TRAP  = uint32(0x00030000) /* disallow and force a SIGSYS */
	SECCOMP_RET_ERRNO = uint32(0x00050000) /* returns an errno */
	SECCOMP_RET_TRACE = uint32(0x7ff00000) /* pass to a tracer or disallow */
	SECCOMP_RET_ALLOW = uint32(0x7fff0000) /* allow */
)
