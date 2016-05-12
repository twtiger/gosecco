package compiler2

import (
	"strconv"
	"strings"

	"github.com/twtiger/gosecco/constants"
)

// actionDescriptionToK turns string specifications of return actions into compiled values acceptable for the compiler to insert
func actionDescriptionToK(v string) uint32 {
	switch strings.ToLower(v) {
	case "trap":
		return SECCOMP_RET_TRAP
	case "kill":
		return SECCOMP_RET_KILL
	case "allow":
		return SECCOMP_RET_ALLOW
	case "trace":
		return SECCOMP_RET_TRACE
	}

	if res, err := strconv.ParseUint(v, 0, 16); err == nil {
		return SECCOMP_RET_ERRNO | uint32(res)
	}

	if res, ok := constants.GetError(v); ok {
		return res
	}

	return 0
}
