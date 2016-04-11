package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	 "strconv"
)

var tokenTypes = make(map[token.Token]string)

func init() {
	tokenTypes[token.LOR] = "booleanArguments"
	tokenTypes[token.LAND] = "booleanArguments"

	tokenTypes[token.GTR] = "integerArguments"
	tokenTypes[token.EQL] = "integerArguments"
}


func surround(s string) string {
	return "func() { " + s + "}"
}


func parseRule(s string) (rule, error) {
	p := strings.SplitN(s, ":", 2)
	name, rl := p[0], p[1]
	fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(surround(rl))
	ast.Print(fs, tr)
	return rule{name, unwrapToplevel(tr)}, nil
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
    		return argumentNode{0}
    	}
    case *ast.BasicLit:
    	i, _ := strconv.Atoi(f.Value)
    	return literalNode{i}
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
			}
		} else if takesIntegerArguments(f) {
			var cmp string;
			switch f.Op {
				case token.GTR: cmp = "gt"
				case token.EQL: cmp = "eq"
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
