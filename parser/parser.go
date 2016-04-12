package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/twtiger/go-seccomp/tree"
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
		tokenTypes[k] = "integerArguments"
	}
}

func surround(s string) string {
	return "func() { " + s + "}"
}

// func parseRule(s string) (rule, error) {
// 	p := strings.SplitN(s, ":", 2)
// 	name, expr := p[0], p[1]
// 	e, _ := parseExpression(expr)
// 	return rule{name, e}, nil
// }

func parseExpression(expr string) (tree.Expression, error) {
	// fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(surround(expr))
	// ast.Print(fs, tr)
	return unwrapToplevel(tr), nil
}

func unwrapToplevel(x ast.Node) tree.Expression {
	switch f := x.(type) {
	case *ast.FuncLit:
		return unwrapToplevel(f.Body)
	case *ast.BlockStmt:
		return unwrapToplevel(f.List[0])
	case *ast.ExprStmt:
		return unwrapBooleanExpression(f.X)
	default:
		panic("Not a valid top level statement")
	}
}

func unwrapIntegerExpression(x ast.Node) tree.Numeric {
	switch f := x.(type) {
	case *ast.Ident:
		switch f.Name {
		case "arg0":
			return tree.Argument{0}
		case "arg1":
			return tree.Argument{1}
		}
		// TODO: More arguments here
		// TODO: variables possible here
	case *ast.BasicLit:
		// TODO: errors here
		i, _ := strconv.Atoi(f.Value)
		return tree.NumericLiteral{uint32(i)}
	case *ast.BinaryExpr:
		left := unwrapIntegerExpression(f.X)
		right := unwrapIntegerExpression(f.Y)
		op := arithmeticOps[f.Op]
		// TODO: handle operators we don't support here

		return tree.Arithmetic{Left: left, Right: right, Op: op}
	default:
		panic("No integer")
	}
	panic("Not a valid integer expression")
}

func takesBooleanArguments(f *ast.BinaryExpr) bool {
	return tokenTypes[f.Op] == "booleanArguments"
}

func takesIntegerArguments(f *ast.BinaryExpr) bool {
	return tokenTypes[f.Op] == "integerArguments"
}

func unwrapBooleanExpression(x ast.Node) tree.Boolean {
	switch f := x.(type) {
	case *ast.BasicLit:
		switch f.Value {
		// TODO: Handle other values here
		case "1":
			return tree.BooleanLiteral{true}
		}
		// TODO: handle failure here
	case *ast.BinaryExpr:
		if takesBooleanArguments(f) {
			switch f.Op {
			case token.LOR:
				left := unwrapBooleanExpression(f.X)
				right := unwrapBooleanExpression(f.Y)
				return tree.Or{Left: left, Right: right}
			case token.LAND:
				left := unwrapBooleanExpression(f.X)
				right := unwrapBooleanExpression(f.Y)
				return tree.And{Left: left, Right: right}
			}
		} else if takesIntegerArguments(f) {
			cmp := comparisonOps[f.Op]
			// TODO: handle incorrect thingy here
			left := unwrapIntegerExpression(f.X)
			right := unwrapIntegerExpression(f.Y)
			return tree.Comparison{Left: left, Right: right, Op: cmp}
		}
	default:
		panic(fmt.Sprintf("can't do this with %#v", x))
	}
	panic("Not a valid boolean expression")
}
