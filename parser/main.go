package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

var tokenTypes = make(map[token.Token]string)
var arithmeticOps = make(map[token.Token]string)

func buildArithmeticOps() {
	arithmeticOps[token.ADD] = "add"
	arithmeticOps[token.SUB] = "sub"
	arithmeticOps[token.MUL] = "mul"
	arithmeticOps[token.QUO] = "quo"
	arithmeticOps[token.REM] = "rem"
	arithmeticOps[token.AND] = "and"
	arithmeticOps[token.OR] = "or"
	arithmeticOps[token.XOR] = "xor"
	arithmeticOps[token.SHL] = "shl"
	arithmeticOps[token.SHR] = "shr"
}

func init() {
	tokenTypes[token.LOR] = "booleanArguments"
	tokenTypes[token.LAND] = "booleanArguments"

	tokenTypes[token.GTR] = "integerArguments"
	tokenTypes[token.EQL] = "integerArguments"

	buildArithmeticOps()
}

func surround(s string) string {
	return "func() { " + s + "}"
}

func parseRule(s string) (rule, error) {
	p := strings.SplitN(s, ":", 2)
	name, expr := p[0], p[1]
	e, _ := parseExpression(expr)
	return rule{name, e}, nil
}

func parseExpression(expr string) (expression, error) {
	fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(surround(expr))
	ast.Print(fs, tr)
	return unwrapToplevel(tr), nil
}

func unwrapToplevel(x ast.Node) expression {
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

func unwrapIntegerExpression(x ast.Node) integerExpression {
	switch f := x.(type) {
	case *ast.Ident:
		switch f.Name {
		case "arg0":
			return argumentNode{index: 0}
		case "arg1":
			return argumentNode{index: 1}
		}
	case *ast.BasicLit:
		i, _ := strconv.Atoi(f.Value)
		return literalNode{i}
	case *ast.BinaryExpr:
		left := unwrapIntegerExpression(f.X)
		right := unwrapIntegerExpression(f.Y)
		op, _ := arithmeticOps[f.Op]
		return arithmetic{left, right, op}
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

func unwrapBooleanExpression(x ast.Node) booleanExpression {
	switch f := x.(type) {
	case *ast.BasicLit:
		switch f.Value {
		case "1":
			return trueLiteral{}
		}
	case *ast.BinaryExpr:
		if takesBooleanArguments(f) {
			switch f.Op {
			case token.LOR:
				left := unwrapBooleanExpression(f.X)
				right := unwrapBooleanExpression(f.Y)
				return orExpr{left, right}
			case token.LAND:
				left := unwrapBooleanExpression(f.X)
				right := unwrapBooleanExpression(f.Y)
				return andExpr{left, right}
			}
		} else if takesIntegerArguments(f) {
			var cmp string
			switch f.Op {
			case token.GTR:
				cmp = "gt"
			case token.EQL:
				cmp = "eq"
			}
			left := unwrapIntegerExpression(f.X)
			right := unwrapIntegerExpression(f.Y)
			return comparison{left, right, cmp}

		}
	default:
		panic(fmt.Sprintf("can't do this with %#v", x))
	}
	panic("Not a valid boolean expression")
}
