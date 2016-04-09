package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"fmt"
)

func surround(s string) string {
	return "func() { " + s + "}"
}

func main() {
	fs := token.NewFileSet()
	tr, _ := parser.ParseExpr(surround("arg0 == 42"))
	ast.Print(fs, tr)

	parseRule(surround("arg0 != 42"))
}

func parseRule(s string) BooleanExpression {
	tr, _ := parser.ParseExpr(surround(s))
	return translateBooleanExpression(tr)
}

func unwrapIgnored(x ast.Node) ast.Node {
	switch f := x.(type) {
	case *ast.FuncLit:
		return unwrapIgnored(f.Body)
	case *ast.BlockStmt:
		return unwrapIgnored(f.List[0])
	case *ast.ExprStmt:
		return unwrapIgnored(f.X)
	default:
		return x
	}
}

func translateBooleanExpression(x ast.Node) BooleanExpression {
	x = unwrapIgnored(x)

	switch f := x.(type) {
	case *ast.BinaryExpr:
		switch f.Op {
		case token.EQL:
			fmt.Printf("WOOT EQUALS\n")
		}
	}
	return nil
}
