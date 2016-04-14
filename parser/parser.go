package parser

import (
	"errors"
	"fmt"
	"github.com/twtiger/go-seccomp/tree"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
)

var tokenTypes = make(map[token.Token]string)
var arithmeticOps = make(map[token.Token]tree.ArithmeticType)
var comparisonOps = make(map[token.Token]tree.ComparisonType)

func buildArithmeticOps() {
	arithmeticOps[token.ADD] = tree.PLUS
	arithmeticOps[token.SUB] = tree.MINUS
	arithmeticOps[token.MUL] = tree.MULT
	arithmeticOps[token.QUO] = tree.DIV
	arithmeticOps[token.REM] = tree.MOD
	arithmeticOps[token.AND] = tree.BINAND
	arithmeticOps[token.OR] = tree.BINOR
	arithmeticOps[token.XOR] = tree.BINXOR
	arithmeticOps[token.SHL] = tree.LSH
	arithmeticOps[token.SHR] = tree.RSH
}

func buildComparisonOps() {
	comparisonOps[token.EQL] = tree.EQL
	comparisonOps[token.LSS] = tree.LT
	comparisonOps[token.GTR] = tree.GT
	comparisonOps[token.NEQ] = tree.NEQL
	comparisonOps[token.LEQ] = tree.LTE
	comparisonOps[token.GEQ] = tree.GTE
	comparisonOps[token.AND] = tree.BIT
}

func init() {
	buildArithmeticOps()
	buildComparisonOps()

	tokenTypes[token.LOR] = "booleanArguments"
	tokenTypes[token.LAND] = "booleanArguments"

	for k := range arithmeticOps {
		tokenTypes[k] = "integerArguments"
	}

	for k := range comparisonOps {
		tokenTypes[k] = "numericArguments"
	}
}

func surround(s string) string {
	return "func() { " + s + "}"
}

func parseExpression(expr string) (tree.Expression, error) {
	// fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(surround(expr))
	// ast.Print(fs, tr)
	parsedtree, err := unwrapToplevel(tr)
	return parsedtree, err
}

func unwrapToplevel(x ast.Node) (tree.Expression, error) {
	switch f := x.(type) {
	case *ast.FuncLit:
		return unwrapToplevel(f.Body)
	case *ast.BlockStmt:
		return unwrapToplevel(f.List[0])
	case *ast.ExprStmt:
		return unwrapBooleanExpression(f.X)
	default:
		// panicWithInfo(x)
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}

var argRegexpRE = regexp.MustCompile(`^arg([0-5])$`)

func identExpression(f *ast.Ident) (tree.Numeric, error) {
	if match := argRegexpRE.FindStringSubmatch(f.Name); match != nil {
		ix, _ := strconv.Atoi(match[1])
		return tree.Argument{ix}, nil
	}
	switch f.Name {
	case "true":
		return tree.BooleanLiteral{true}, nil
	case "false":
		return tree.BooleanLiteral{false}, nil
	// Handle other cases here
	default:
		return tree.Variable{f.Name}, nil
	}
	return tree.Variable{f.Name}, nil
}

func unwrapNumericExpression(x ast.Node) (tree.Numeric, error) {
	switch f := x.(type) {
	case *ast.Ident:
		// Ensure ident doesn't contain stupidness like packages and stuff
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
		// TODO: Fail in a good way here
	default:
		// panicWithInfo(x)
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

func inclusionExpression(f *ast.CallExpr) (tree.Boolean, error) {
	var pos bool
	var left tree.Numeric
	var err error
	var right[] tree.Numeric
	var val tree.Numeric

	switch p := f.Fun.(type) {
	case *ast.Ident:
		if p.Name == "in" {
			pos = true
		}
		if p.Name == "notIn" {
			pos = false
		}
	}

	switch p := f.Args[0].(type) {
	case *ast.Ident:
		left, err = identExpression(p)
	case *ast.BasicLit:
		left, err = unwrapNumericExpression(p)
	}

	for _, e := range f.Args {
		switch y := e.(type) {
		case *ast.BasicLit:
			val, err = unwrapNumericExpression(y)
			right = append(right, val)
		}
	}

	if err != nil {
		return nil, err
	}

	return tree.Inclusion{pos, left, right}, nil
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
		return inclusionExpression(f)
	case *ast.Ident:
		return identExpression(f)
		// TODO: Fail in a good way here
	default:
		// panicWithInfo(x)
	}
	return nil, errors.New("Expression is invalid. Unable to parse.")
}
