package parser

import . "gopkg.in/check.v1"

type RuleSuite struct{}

var _ = Suite(&RuleSuite{})

func parseRuleHeadCheck(c *C, s string, r ruleHead) {
	res, ok := parseRuleHead(s)
	c.Assert(ok, Equals, true)
	c.Check(res, Equals, r)
}

func (s *RuleSuite) Test_parseRuleHead_parsesValidRuleHeads(c *C) {
	parseRuleHeadCheck(c, "read", ruleHead{syscall: "read"})
	parseRuleHeadCheck(c, "write", ruleHead{syscall: "write"})
	parseRuleHeadCheck(c, "\t write  ", ruleHead{syscall: "write"})
	parseRuleHeadCheck(c, "fcntl[]", ruleHead{syscall: "fcntl"})
	parseRuleHeadCheck(c, "fcntl [ ] ", ruleHead{syscall: "fcntl"})
	parseRuleHeadCheck(c, " fcntl [ +kill ] ", ruleHead{syscall: "fcntl", positive: "kill"})
	parseRuleHeadCheck(c, " fcntl[ -kill] ", ruleHead{syscall: "fcntl", negative: "kill"})
	parseRuleHeadCheck(c, " fcntl[ -kill, +trace] ", ruleHead{syscall: "fcntl", negative: "kill", positive: "trace"})
	parseRuleHeadCheck(c, " fcntl[+trace,-kill] ", ruleHead{syscall: "fcntl", negative: "kill", positive: "trace"})
	parseRuleHeadCheck(c, " fcntl[+trace,-42] ", ruleHead{syscall: "fcntl", negative: "42", positive: "trace"})
}
