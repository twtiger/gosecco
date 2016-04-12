package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
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
		tokenTypes[k] = "numericArguments"
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
		panicWithInfo(x)
	}
	return nil
}

var argRegexpRE = regexp.MustCompile(`^arg([0-5])$`)

func identExpression(f *ast.Ident) tree.Numeric {
	if match := argRegexpRE.FindStringSubmatch(f.Name); match != nil {
		ix, _ := strconv.Atoi(match[1])
		return tree.Argument{ix}
	}
	return tree.Variable{f.Name}
}

func unwrapNumericExpression(x ast.Node) tree.Numeric {
	switch f := x.(type) {
	case *ast.Ident:
		// Ensure ident doesn't contain stupidness like packages and stuff
		return identExpression(f)
	case *ast.BasicLit:
		// TODO: errors here
		i, _ := strconv.Atoi(f.Value)
		return tree.NumericLiteral{uint32(i)}
	case *ast.BinaryExpr:
		left := unwrapNumericExpression(f.X)
		right := unwrapNumericExpression(f.Y)
		op := arithmeticOps[f.Op]
		// TODO: handle operators we don't support here
		return tree.Arithmetic{Left: left, Right: right, Op: op}
	case *ast.ParenExpr:
		return unwrapNumericExpression(f.X)
	case *ast.UnaryExpr:
		operand := unwrapNumericExpression(f.X)
		if f.Op == token.XOR {
			return tree.BinaryNegation{operand}
		}
		// TODO: Fail in a good way here
	default:
		panicWithInfo(x)
	}
	return nil
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

func booleanArgExpression(f *ast.BinaryExpr) tree.Boolean {
	left := unwrapBooleanExpression(f.X)
	right := unwrapBooleanExpression(f.Y)
	switch f.Op {
	case token.LOR:
		return tree.Or{Left: left, Right: right}
	case token.LAND:
		return tree.And{Left: left, Right: right}
	}
	panic("Not recognized operator for boolean arg expression. This shouldn't be possible")
}

func numericArgExpression(f *ast.BinaryExpr) tree.Boolean {
	cmp := comparisonOps[f.Op]
	// TODO: handle incorrect thingy here
	left := unwrapNumericExpression(f.X)
	right := unwrapNumericExpression(f.Y)
	return tree.Comparison{Left: left, Right: right, Op: cmp}
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
			return booleanArgExpression(f)
		} else if takesNumericArguments(f) {
			return numericArgExpression(f)
		}
	case *ast.ParenExpr:
		return unwrapBooleanExpression(f.X)
	case *ast.UnaryExpr:
		operand := unwrapBooleanExpression(f.X)
		if f.Op == token.NOT {
			return tree.Negation{operand}
		}
		// TODO: Fail in a good way here
	default:
		panicWithInfo(x)
	}
	return nil
}
