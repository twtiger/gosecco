package parser2

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TokenizerSuite struct{}

var _ = Suite(&TokenizerSuite{})

func (s *TokenizerSuite) Test_bla1(c *C) {
	//	var result string
	parse([]byte("0 1 2 3 42 32432 0B1 07 0xFFFFFFAc3 [](), , arg0 arg9 TRUE | & flalse != == & ! || && + << ~ ^//& inx in notIn"),
		func(tt token, td []byte) {
			fmt.Printf("Token: %s data: %s\n", tokens[tt], string(td))
		})

	//	c.Assert(result, Equals, "INTDEC")
}
