package compiler

const LOAD = BPF_LD | BPF_W | BPF_ABS
const LOAD_VAL = BPF_LD | BPF_IMM

const JEQ_K = BPF_JMP | BPF_JEQ | BPF_K
const JEQ_X = BPF_JMP | BPF_JEQ | BPF_X

const JEG_K = BPF_JMP | BPF_JGT | BPF_K
const JEG_X = BPF_JMP | BPF_JGT | BPF_X

const JEGE_K = BPF_JMP | BPF_JGE | BPF_K

const JEGE_X = BPF_JMP | BPF_JGE | BPF_X

const JSET_K = BPF_JMP | BPF_JSET | BPF_K
const JSET_X = BPF_JMP | BPF_JSET | BPF_X

const ADD_K = BPF_ALU | BPF_ADD | BPF_K
const SUB_K = BPF_ALU | BPF_SUB | BPF_K
const MUL_K = BPF_ALU | BPF_MUL | BPF_K
const DIV_K = BPF_ALU | BPF_DIV | BPF_K
const AND_K = BPF_ALU | BPF_AND | BPF_K
const OR_K = BPF_ALU | BPF_OR | BPF_K

const LSH_K = BPF_ALU | BPF_LSH | BPF_K
const RSH_K = BPF_ALU | BPF_RSH | BPF_K

// const MOD_K = BPF_ALU | BPF_MOD | BPF_K
// const XOR_K = BPF_ALU | BPF_XOR | BPF_K

const RET_K = BPF_RET | BPF_K
const A_TO_X = BPF_MISC | BPF_TAX
