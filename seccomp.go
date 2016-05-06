package gosecco

import (
	"fmt"

	"github.com/twtiger/gosecco/data"
	"github.com/twtiger/gosecco/native"

	"golang.org/x/sys/unix"
)

// CheckSupport checks for the required seccomp support in the kernel.
func CheckSupport() error {
	// TODO: no testing really possible
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
	// TODO: test when compiler is ready:
	// - test that parser errors come through
	// - test that unification works and that errors come through
	// - test that default pos and neg actions come through
	// - test that the type checker errors come through
	// - test that the simplifier is invoked and simplifies stuff
	// - test that simplifier errors come through
	// - test that the compiler works and returns the expected results
	// - test that compiler errors come through
	return nil, nil
}

// Compile provides the compatibility interface for gosecco - it has the same signature as
// Compile from the go-seccomp package and should provide the same behavior.
// However, the modern interface is through the Prepare function
func Compile(path string, enforce bool) ([]unix.SockFilter, error) {
	// TODO: test once compiler is done, light testing needed, since main testing
	// will be of the Prepare method
	return nil, nil
}

// CompileBlacklist provides the compatibility interface for gosecco, for blacklist mode
// It has the same signature as CompileBlacklist from Subgraphs go-seccomp and should provide the same behavior.
// However, the modern interface is through the Prepare function
func CompileBlacklist(path string, enforce bool) ([]unix.SockFilter, error) {
	// TODO: test once compiler is done, light testing needed, since main testing
	// will be of the Prepare method
	return nil, nil
}

// Load makes the seccomp system call to install the bpf filter for
// all threads (with tsync). Most users of this library should use
// Install instead of Load, since Install ensures that prctl(set_no_new_privs, 1)
// has been called
func Load(bpf []unix.SockFilter) error {
	if size, limit := len(bpf), 0xffff; size > limit {
		return fmt.Errorf("filter program too big: %d bpf instructions (limit = %d)", size, limit)
	}
	prog := &data.SockFprog{
		Filter: &bpf[0],
		Len:    uint16(len(bpf)),
	}

	return native.InstallSeccomp(prog)
}

// Install will install the given policy filters into the kernel
func Install(bpf []unix.SockFilter) error {
	// TODO, doesn't need testing (not really possible)
	return nil
}

// InstallBlacklist makes the necessary system calls to install the Seccomp-BPF
// filter for the current process (all threads). Install can be called
// multiple times to install additional filters.
func InstallBlacklist(bpf []unix.SockFilter) error {
	// TODO, doesn't need testing (not really possible)
	return nil
}
