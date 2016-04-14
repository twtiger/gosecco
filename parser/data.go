package parser

import (
	"go/token"

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
