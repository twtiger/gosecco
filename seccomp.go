package gosecco

import "golang.org/x/sys/unix"

// CheckSupport checks for the required seccomp support in the kernel.
func CheckSupport() error {
	return nil
}

// SeccompSettings contains the extra settings necessary to tweak the
// behavior of the compilation process
type SeccompSettings struct {
	extraDefinitions      []string
	defaultPositiveAction string
	defaultNegativeAction string
}

// Prepare will take the given path and settings, parse and compile the given
// data, combined with the settings - and returns the bytecode
func Prepare(path string, s SeccompSettings) ([]unix.SockFilter, error) {
	return nil, nil
}

// Compile provides the compatibility interface for gosecco - it has the same signature as
// Compile from the go-seccomp package and should provide the same behavior.
// However, the modern interface is through the Prepare function
func Compile(path string, enforce bool) ([]unix.SockFilter, error) {
	return nil, nil
}

// Load makes the seccomp system call to install the bpf filter for
// all threads (with tsync). Most users of this library should use
// Install instead of Load, since Install ensures that prctl(set_no_new_privs, 1)
// has been called
func Load(bpf []unix.SockFilter) error {
	return nil
}

// Install will install the given policy filters into the kernel
func Install(bpf []unix.SockFilter) error {
	return nil
}
