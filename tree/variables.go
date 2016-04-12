package tree

// Argument represents an argment given to the syscall
type Argument struct {
	Index int
}

// Accept implements Expression
func (v Argument) Accept(vs Visitor) {
	vs.AcceptArgument(v)
}

// Variable represents a variable used before
type Variable struct {
	Name string
}

// Accept implements Expression
func (v Variable) Accept(vs Visitor) {
	vs.AcceptVariable(v)
}
