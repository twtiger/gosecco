package gosecco

import (
	"os"
	"path"
	"strings"
	"testing"

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
