package data

type SeccompData struct {
	NR                 int32     // The system call number.
	Arch               uint32    // System call convention as an AUDIT_ARCH_* value.
	InstructionPointer uint64    // At the time of the system call.
	Args               [6]uint64 // System call arguments (always stored as 64-bit values).
}

type SockFilter struct {
	Code uint16 // Actual filter code.
	JT   uint8  // Jump true.
	JF   uint8  // Jump false.
	K    uint32 // Generic multiuse field.
}
