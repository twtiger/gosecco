package gosecco

import (
	"fmt"
	"runtime"
	"syscall"

	"github.com/twtiger/gosecco/checker"
	"github.com/twtiger/gosecco/compiler"
	"github.com/twtiger/gosecco/data"
	"github.com/twtiger/gosecco/native"
	"github.com/twtiger/gosecco/parser"
	"github.com/twtiger/gosecco/precompilation"
	"github.com/twtiger/gosecco/simplifier"
	"github.com/twtiger/gosecco/tree"
	"github.com/twtiger/gosecco/unifier"

	"golang.org/x/sys/unix"
)

// CheckSupport checks for the required seccomp support in the kernel.
func CheckSupport() error {
	if err := native.CheckGetSeccomp(); err != nil {
		return fmt.Errorf("seccomp not available: %v", err)
	}
	if err := native.CheckSetSeccompModeFilter(); err != syscall.EFAULT {
		return fmt.Errorf("seccomp filter not available: %v", err)
	}
	if err := native.CheckSetSeccompModeFilterWithSeccomp(); err != syscall.EFAULT {
		return fmt.Errorf("seccomp syscall not available: %v", err)
	}
	if err := native.CheckSetSeccompModeTsync(); err != syscall.EFAULT {
		return fmt.Errorf("seccomp tsync not available: %v", err)
	}
	return nil
}

// SeccompSettings contains the extra settings necessary to tweak the
// behavior of the compilation process
type SeccompSettings struct {
	// ExtraDefinitions contains paths to files with extra definitions to parse
	ExtraDefinitions      []string
	DefaultPositiveAction string
	DefaultNegativeAction string
	DefaultPolicyAction   string
	ActionOnX32           string
	ActionOnAuditFailure  string
}

// Prepare will take the given path and settings, parse and compile the given
// data, combined with the settings - and returns the bytecode
func Prepare(path string, s SeccompSettings) ([]unix.SockFilter, error) {
	// TODO: test when compiler is ready:
	// - test that default pos and neg actions come through
	// - test that the type checker errors come through
	// - test that the simplifier is invoked and simplifies stuff
	// - test that simplifier errors come through
	// - test panic when default settings are not set
	// - test that the compiler works and returns the expected results
	// - test that compiler errors come through

	// Parsing of extra files with definitions
	extras := make([]map[string]tree.Macro, len(s.ExtraDefinitions))
	for ix, ed := range s.ExtraDefinitions {
		rp, e := parser.ParseFile(ed)
		if e != nil {
			return nil, e
		}
		p, e2 := unifier.Unify(rp, nil, "", "", "")
		if e2 != nil {
			return nil, e2
		}
		extras[ix] = p.Macros
	}

	// Parsing
	rp, err := parser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Unifying
	pol, err := unifier.Unify(rp, extras, s.DefaultPositiveAction, s.DefaultNegativeAction, s.DefaultPolicyAction)
	if err != nil {
		return nil, err
	}

	// Type checking
	errors := checker.EnsureValid(pol)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	// Simplification
	simplifier.SimplifyPolicy(pol)

	// Pre-compilation
	errors = precompilation.EnsureValid(pol)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	// Compilation
	return compiler.Compile(pol)
}

// Compile provides the compatibility interface for gosecco - it has the same signature as
// Compile from the go-seccomp package and should provide the same behavior.
// However, the modern interface is through the Prepare function
func Compile(path string, enforce bool) ([]unix.SockFilter, error) {
	// TODO: test when compiler is done
	// TODO: set all three default actions correctly

	settings := SeccompSettings{}
	settings.DefaultPositiveAction = "allow"
	settings.ActionOnAuditFailure = "kill"
	if enforce {
		settings.DefaultNegativeAction = "kill"
		settings.DefaultPolicyAction = "kill"
	} else {
		settings.DefaultNegativeAction = "trace"
		settings.DefaultPolicyAction = "trace"
	}

	return Prepare(path, settings)
}

// CompileBlacklist provides the compatibility interface for gosecco, for blacklist mode
// It has the same signature as CompileBlacklist from Subgraphs go-seccomp and should provide the same behavior.
// However, the modern interface is through the Prepare function
func CompileBlacklist(path string, enforce bool) ([]unix.SockFilter, error) {
	// TODO: test when compiler is done
	// TODO: set all three default actions correctly

	settings := SeccompSettings{}
	settings.DefaultNegativeAction = "allow"
	settings.DefaultPolicyAction = "allow"
	settings.ActionOnX32 = "kill"
	settings.ActionOnAuditFailure = "kill"
	if enforce {
		settings.DefaultPositiveAction = "kill"
	} else {
		settings.DefaultPositiveAction = "trace"
	}

	return Prepare(path, settings)
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
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := native.NoNewPrivs(); err != nil {
		return err
	}
	return Load(bpf)
}

// InstallBlacklist makes the necessary system calls to install the Seccomp-BPF
// filter for the current process (all threads). Install can be called
// multiple times to install additional filters.
func InstallBlacklist(bpf []unix.SockFilter) error {
	return Install(bpf)
}
