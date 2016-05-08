package compiler2

import "golang.org/x/sys/unix"

type compilerContext struct {
	result []unix.SockFilter
}

func (c *compilerContext) add(a []unix.SockFilter) {
	c.result = append(c.result, a...)
}

func (c *compilerContext) push(a unix.SockFilter) {
	c.result = append(c.result, a)
}

func (c *compilerContext) op(code uint16, k uint32) {
	c.push(unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
}
