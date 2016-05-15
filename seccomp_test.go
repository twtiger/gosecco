package gosecco

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/twtiger/gosecco/asm"
	"golang.org/x/sys/unix"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SeccompSuite struct{}

var _ = Suite(&SeccompSuite{})

func (s *SeccompSuite) Test_loadingTooBigBpf(c *C) {
	inp := make([]unix.SockFilter, 0xFFFF+1)
	res := Load(inp)
	c.Assert(res, ErrorMatches, "filter program too big: 65536 bpf instructions \\(limit = 65535\\)")
}

func getActualTestFolder() string {
	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, "/parser") {
		return wd
	}
	return path.Join(wd, "parser")
}

func (s *SeccompSuite) Test_parseInvalidFileReturnsErrors(c *C) {
	set := SeccompSettings{}
	f := getActualTestFolder() + "/failing_test_policy"
	_, ee := Prepare(f, set)
	c.Assert(ee, ErrorMatches, ".*parser/failing_test_policy:1: unexpected end of line")
}

func (s *SeccompSuite) Test_parseUnificationErrorReturnsError(c *C) {
	set := SeccompSettings{}
	f := getActualTestFolder() + "/missing_variable_policy"
	_, ee := Prepare(f, set)
	c.Assert(ee, ErrorMatches, "Variable not defined")
}

func (s *SeccompSuite) Test_parseValidPolicyFile(c *C) {
	set := SeccompSettings{DefaultPositiveAction: "allow", DefaultNegativeAction: "kill", DefaultPolicyAction: "kill"}
	f := getActualTestFolder() + "/valid_test_policy"
	res, ee := Prepare(f, set)

	c.Assert(ee, Equals, nil)

	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs\t4\n"+
		"jeq_k\t00\t05\tC000003E\n"+
		"ld_abs\t0\n"+
		"jeq_k\t00\t01\t1\n"+
		"jmp\t1\n"+
		"jmp\t1\n"+
		"ret_k\t7FFF0000\n"+
		"ret_k\t0\n")
}
