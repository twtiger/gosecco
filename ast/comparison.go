package ast

// ComparisonType specifies the possible comparison types
type ComparisonType int

// Contains all the comparison types
const (
	EQL ComparisonType = iota
	NEQL
	GT
	GTE
	LT
	LTE
	BIT
)

// ComparisonNames maps types to names for presentation
var ComparisonNames = map[ComparisonType]string{
	EQL:  "==",
	NEQL: "!=",
	GT:   ">",
	GTE:  ">=",
	LT:   "<",
	LTE:  "<=",
	BIT:  "&",
}

// Comparison represents a comparison
type Comparison struct {
	Op          ComparisonType
	Left, Right Numeric
}

// Accept implements Expression
func (v Comparison) Accept(vs Visitor) {
	vs.AcceptComparison(v)
}
