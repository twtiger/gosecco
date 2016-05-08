package compiler2

import "syscall"

const OP_LOAD_VAL = syscall.BPF_LD | syscall.BPF_IMM
