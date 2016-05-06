package native

import (
	"syscall"
	"unsafe"

	"github.com/twtiger/gosecco/data"
)

// #include <linux/seccomp.h>
import "C"

// seccomp is a wrapper for the 'seccomp' system call.
// See <linux/seccomp.h> for valid op and flag values.
// uargs is typically a pointer to struct sock_fprog.
func seccomp(op, flags uintptr, uargs unsafe.Pointer) error {
	_, _, e := syscall.Syscall(syscall.PR_GET_SECCOMP, op, flags, uintptr(uargs))
	if e != 0 {
		return e
	}
	return nil
}

// InstallSeccomp using native methods
func InstallSeccomp(prog *data.SockFprog) error {
	return seccomp(C.SECCOMP_SET_MODE_FILTER, C.SECCOMP_FILTER_FLAG_TSYNC, unsafe.Pointer(prog))
}
