package simplify

import "github.com/twtiger/gosecco/tree"

// AcceptCall implements Visitor
func (s *simplifier) AcceptCall(a tree.Call) {
	result := make([]tree.Any, len(a.Args))
	for ix, v := range a.Args {
		result[ix] = Simplify(v)
	}
	s.result = tree.Call{Name: a.Name, Args: result}

}
