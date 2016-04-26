package tree

// RawPolicy represents the raw parsed rules and macros in the order they were encountered. This can be used to generate the final Policy
type RawPolicy struct {
	ListType     PolicyType
	RuleOrMacros []interface{}
}

type PolicyType int

const (
	WhiteList PolicyType = iota
	BlackList
)

// Policy represents a complete policy file. It is possible to combine more than one policy file
type Policy struct {
	DefaultPositiveAction string
	DefaultNegativeAction string
	Macros                map[string]Macro
	Rules                 []Rule
}
