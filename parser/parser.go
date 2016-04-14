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

func parseExpression(expr string) (tree.Expression, bool, uint16, error) {
	// fset := token.NewFileSet()
	tr, err := parser.ParseExpr(surround(expr))
	if err != nil {
		return nil, false, 0, errors.New("Expression is invalid. Unable to parse.")
	}
	// ast.Print(fset, tr)
	return unwrapToplevel(tr)
}

func extractReturnInformation(f *ast.ReturnStmt) (bool, uint16) {
	if len(f.Results) > 0 {
		val, ok := f.Results[0].(*ast.BasicLit)
		if ok && val.Kind == token.INT {
			errno, err := strconv.ParseUint(val.Value, 0, 16)
			if err == nil {
				return true, uint16(errno)
			}
		}
	}
	return false, 0
}

func unwrapToplevel(x ast.Node) (tree.Expression, bool, uint16, error) {
	body := x.(*ast.FuncLit).Body.List

	var res tree.Expression
	var err error

	switch f := body[0].(type) {
	case *ast.ExprStmt:
		res, err = unwrapBooleanExpression(f.X)
		if err != nil {
			return nil, false, 0, err
		}
	case *ast.ReturnStmt:
		b, v := extractReturnInformation(f)
		if b {
			return nil, b, v, nil
		}
	}

	if len(body) > 1 {
		f2, ok := body[1].(*ast.ReturnStmt)
		if ok {
			b, v := extractReturnInformation(f2)
			if b {
				return res, b, v, nil
			}
		}
	}

	if res != nil {
		return res, false, 0, nil
	}

	return nil, false, 0, errors.New("Expression is invalid. Unable to parse.")
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
		if f.Kind == token.INT {
			i, _ := strconv.ParseUint(f.Value, 0, 32)
			return tree.NumericLiteral{uint32(i)}, nil
		}
		return nil, errors.New("Invalid literal type - this language only supports numbers")
	case *ast.BinaryExpr:
		op, ok := arithmeticOps[f.Op]
		if !ok {
			return nil, fmt.Errorf("Operator '%s' cannot be used in a numeric context", f.Op)
		}

		left, err := unwrapNumericExpression(f.X)
		if err != nil {
			return nil, err
		}

		right, err := unwrapNumericExpression(f.Y)
		if err != nil {
			return nil, err
		}

		return tree.Arithmetic{Left: left, Right: right, Op: op}, nil
	case *ast.ParenExpr:
		return unwrapNumericExpression(f.X)
	case *ast.CallExpr:
		return callExpression(f)
	case *ast.UnaryExpr:
		operand, err := unwrapNumericExpression(f.X)
		if err != nil {
			return nil, err
		}
		if f.Op == token.XOR {
			return tree.BinaryNegation{operand}, nil
		}
		return nil, fmt.Errorf("Invalid unary operator: '%s'", f.Op)
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}

func takesBooleanArguments(f *ast.BinaryExpr) bool {
	return tokenTypes[f.Op] == "booleanArguments"
}

func takesNumericArguments(f *ast.BinaryExpr) bool {
	return tokenTypes[f.Op] == "numericArguments"
}

func booleanArgExpression(f *ast.BinaryExpr) (tree.Boolean, error) {
	left, err := unwrapBooleanExpression(f.X)
	if err != nil {
		return nil, err
	}

	right, err := unwrapBooleanExpression(f.Y)
	if err != nil {
		return nil, err
	}

	switch f.Op {
	case token.LOR:
		return tree.Or{Left: left, Right: right}, nil
	case token.LAND:
		return tree.And{Left: left, Right: right}, nil
	}
	return nil, fmt.Errorf("Operator '%s' cannot be used in a boolean context", f.Op)
}

func numericArgExpression(f *ast.BinaryExpr) (tree.Boolean, error) {
	cmp, ok := comparisonOps[f.Op]
	if !ok {
		return nil, fmt.Errorf("Operator '%s' cannot be used in a boolean context", f.Op)
	}

	left, err := unwrapNumericExpression(f.X)
	if err != nil {
		return nil, err
	}

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
		case "1":
			return tree.BooleanLiteral{true}, nil
		case "0":
			return tree.BooleanLiteral{false}, nil
		}
		return nil, fmt.Errorf("Invalid boolean literal: '%s'", f.Value)
	case *ast.BinaryExpr:
		if takesBooleanArguments(f) {
			return booleanArgExpression(f)
		} else if takesNumericArguments(f) {
			return numericArgExpression(f)
		} else {
			return nil, fmt.Errorf("Operator '%s' cannot be used in a boolean context", f.Op)
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
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}
