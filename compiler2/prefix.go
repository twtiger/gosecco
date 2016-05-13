package compiler2

func (c *compilerContext) compileAuditArchCheck(on label) {
	// do(bpfLoadArch())
	// do(bpfJeq(auditArch, 1, 0))
	// do(bpfRet(retKill()))
}

func (c *compilerContext) compileX32ABICheck(on label) {
	//  that triggers if the X32_SYSCALL_BIT is set. load NR

	// var X32_SYSCALL_BIT = uint32(C.__X32_SYSCALL_BIT)

	// do(bpfLoadArch())
	// do(bpfJeq(auditArch, 0, 2))

	// do(bpfLoadNR())

	// // Kill if NR > X32_SYSCALL_BIT-1

	// do(bpfJgt(X32_SYSCALL_BIT-1, 0, 1))
	// do(bpfRet(retKill()))
}
