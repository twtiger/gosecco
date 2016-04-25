package tree

type VariableType int

// The different types of argument loads that can happen
const (
	Full VariableType = iota
	Low
	Hi
)
