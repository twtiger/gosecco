package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/twtiger/go-seccomp/tree"
)

func surround(s string) string {
	return "func() { " + s + "}"
}

func parseExpression(expr string) (tree.Expression, error) {
	tr, err := parser.ParseExpr(surround(expr))
	if err != nil {
		return nil, errors.New("Expression is invalid. Unable to parse.")
	}
	return unwrapToplevel(tr)
}

func unwrapToplevel(x ast.Node) (tree.Expression, error) {
	switch f := x.(type) {
	case *ast.FuncLit:
		return unwrapToplevel(f.Body)
	case *ast.BlockStmt:
		return unwrapToplevel(f.List[0])
	case *ast.ExprStmt:
		return unwrapBooleanExpression(f.X)
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}

var argRegexpRE = regexp.MustCompile(`^arg([0-5])$`)

func identExpression(f *ast.Ident) (tree.Numeric, error) {
	if match := argRegexpRE.FindStringSubmatch(f.Name); match != nil {
		// This error _really_ shouldn't be possible
		ix, e := strconv.Atoi(match[1])
		if e != nil {
			panic("Impossible error")
		}
		return tree.Argument{ix}, nil
	}
	switch strings.ToLower(f.Name) {
	case "true":
		return tree.BooleanLiteral{true}, nil
	case "false":
		return tree.BooleanLiteral{false}, nil
	}
	return tree.Variable{f.Name}, nil
}

func unwrapNumericExpression(x ast.Node) (tree.Numeric, error) {
	switch f := x.(type) {
	case *ast.Ident:
		return identExpression(f)
	case *ast.BasicLit:
		// TODO: errors here
		i, _ := strconv.Atoi(f.Value)
		return tree.NumericLiteral{uint32(i)}, nil
	case *ast.BinaryExpr:
		left, err := unwrapNumericExpression(f.X)
		right, err := unwrapNumericExpression(f.Y)
		op := arithmeticOps[f.Op]
		// TODO: handle operators we don't support here
		if err != nil {
			return nil, err
		}
		return tree.Arithmetic{Left: left, Right: right, Op: op}, nil
	case *ast.ParenExpr:
		return unwrapNumericExpression(f.X)
	case *ast.UnaryExpr:
		operand, err := unwrapNumericExpression(f.X)
		if err != nil {
			return nil, err
		}
		if f.Op == token.XOR {
			return tree.BinaryNegation{operand}, nil
		}
		// TODO: other unary expressions possible here?
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}

func panicWithInfo(x interface{}) {
	panic(fmt.Sprintf("sadness: %#v", x))
}

func takesBooleanArguments(f *ast.BinaryExpr) bool {
	return tokenTypes[f.Op] == "booleanArguments"
}

func takesNumericArguments(f *ast.BinaryExpr) bool {
	return tokenTypes[f.Op] == "numericArguments"
}

func booleanArgExpression(f *ast.BinaryExpr) (tree.Boolean, error) {
	left, err := unwrapBooleanExpression(f.X)
	right, err := unwrapBooleanExpression(f.Y)
	switch f.Op {
	case token.LOR:
		return tree.Or{Left: left, Right: right}, nil
	case token.LAND:
		return tree.And{Left: left, Right: right}, nil
	}
	return nil, err
}

func numericArgExpression(f *ast.BinaryExpr) (tree.Boolean, error) {
	cmp := comparisonOps[f.Op]
	// TODO: handle incorrect thingy here
	left, err := unwrapNumericExpression(f.X)
	right, err := unwrapNumericExpression(f.Y)
	if err != nil {
		return nil, err
	}
	return tree.Comparison{Left: left, Right: right, Op: cmp}, nil
}

func expressionsToNumerics(inp []ast.Expr) ([]tree.Numeric, error) {
	var err error
	args := make([]tree.Numeric, len(inp))

	for ix, v := range inp {
		args[ix], err = unwrapNumericExpression(v)
		if err != nil {
			return nil, err
		}
	}

	return args, nil
}

func expressionsToAnys(inp []ast.Expr) ([]tree.Any, error) {
	var err error
	args := make([]tree.Any, len(inp))

	for ix, v := range inp {
		args[ix], err = unwrapNumericExpression(v)
		if err != nil {
			return nil, err
		}
	}

	return args, nil
}

func callExpression(f *ast.CallExpr) (tree.Boolean, error) {
	p, ok := f.Fun.(*ast.Ident)
	if !ok {
		return nil, errors.New("Invalid call expression in boolean context")
	}

	name := p.Name

	switch name {
	case "in", "notIn":
		args, err := expressionsToNumerics(f.Args)
		if err != nil {
			return nil, err
		}
		if len(args) == 0 {
			return nil, errors.New(name + "-expression must have at least a left hand side argument")
		}
		return tree.Inclusion{Positive: name == "in", Left: args[0], Rights: args[1:]}, nil
	default:
		args, err := expressionsToAnys(f.Args)
		if err != nil {
			return nil, err
		}
		return tree.Call{Name: name, Args: args}, nil
	}
}

func unwrapBooleanExpression(x ast.Node) (tree.Boolean, error) {
	switch f := x.(type) {
	case *ast.BasicLit:
		switch f.Value {
		// TODO: Handle other values here
		case "1":
			return tree.BooleanLiteral{true}, nil
		case "0":
			return tree.BooleanLiteral{false}, nil
		}
		// TODO: handle failure here
	case *ast.BinaryExpr:
		if takesBooleanArguments(f) {
			return booleanArgExpression(f)
		} else if takesNumericArguments(f) {
			return numericArgExpression(f)
		}
	case *ast.ParenExpr:
		return unwrapBooleanExpression(f.X)
	case *ast.UnaryExpr:
		operand, err := unwrapBooleanExpression(f.X)
		if err == nil {
			if f.Op == token.NOT {
				return tree.Negation{operand}, nil
			}
		}
	case *ast.CallExpr:
		return callExpression(f)
	case *ast.Ident:
		return identExpression(f)
		// TODO: Fail in a good way here
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}
