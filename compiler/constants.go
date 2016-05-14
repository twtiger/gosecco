package compiler

const (
	SECCOMP_RET_KILL  = uint32(0x00000000) /* kill the task immediately */
	SECCOMP_RET_TRAP  = uint32(0x00030000) /* disallow and force a SIGSYS */
	SECCOMP_RET_ERRNO = uint32(0x00050000) /* returns an errno */
	SECCOMP_RET_TRACE = uint32(0x7ff00000) /* pass to a tracer or disallow */
	SECCOMP_RET_ALLOW = uint32(0x7fff0000) /* allow */
)

const (
	defaultPositive = "allow"
	defaultNegative = "kill"
)

var actionInstructions = map[string]uint32{
	"allow": SECCOMP_RET_ALLOW,
	"kill":  SECCOMP_RET_KILL,
	"trace": SECCOMP_RET_TRACE,
}
