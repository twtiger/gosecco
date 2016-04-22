package parser2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/twtiger/gosecco/tree"
)

// TODO: none of these check errors or whatever

func parseExpression(expr string) tree.Expression {
	tokens, ok := tokenize(expr)
	if !ok {
		// TODO: errors etc
		return nil
	}
	ctx := parseContext{0, tokens, false}
	expression := ctx.logicalORExpression()
	ctx.end()

	return expression
}

func (ctx *parseContext) end() {
	if !ctx.atEnd {
		// Raise error here
	}
}

func (ctx *parseContext) logicalORExpression() tree.Expression {
	left := ctx.logicalANDExpression()
	if ctx.next() == LOR {
		ctx.consume()
		right := ctx.logicalORExpression()
		return tree.Or{Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) logicalANDExpression() tree.Expression {
	left := ctx.inclusiveORExpression()
	if ctx.next() == LAND {
		ctx.consume()
		right := ctx.logicalANDExpression()
		return tree.And{Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) inclusiveORExpression() tree.Expression {
	left := ctx.exclusiveORExpression()
	if ctx.next() == OR {
		ctx.consume()
		right := ctx.inclusiveORExpression()
		return tree.Arithmetic{Op: tree.BINOR, Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) exclusiveORExpression() tree.Expression {
	left := ctx.andExpression()
	if ctx.next() == XOR {
		ctx.consume()
		right := ctx.exclusiveORExpression()
		return tree.Arithmetic{Op: tree.BINXOR, Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) andExpression() tree.Expression {
	left := ctx.equalityExpression()
	if ctx.next() == AND {
		ctx.consume()
		right := ctx.andExpression()
		return tree.Arithmetic{Op: tree.BINAND, Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) equalityExpression() tree.Expression {
	left := ctx.relationalExpression()
	switch ctx.next() {
	case EQL, NEQ:
		op, _ := ctx.consume()
		right := ctx.equalityExpression()
		return tree.Comparison{Op: comparisonOperator[op], Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) relationalExpression() tree.Expression {
	left := ctx.shiftExpression()
	switch ctx.next() {
	case LT, GT, LTE, GTE:
		op, _ := ctx.consume()
		right := ctx.relationalExpression()
		return tree.Comparison{Op: comparisonOperator[op], Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) shiftExpression() tree.Expression {
	left := ctx.additiveExpression()
	switch ctx.next() {
	case LSH, RSH:
		op, _ := ctx.consume()
		right := ctx.shiftExpression()
		return tree.Arithmetic{Op: shiftOperator[op], Left: left, Right: right}
	}

	return left
}

func (ctx *parseContext) additiveExpression() tree.Expression {
	left := ctx.multiplicativeExpression()
	switch ctx.next() {
	case ADD, SUB:
		op, _ := ctx.consume()
		right := ctx.additiveExpression()
		return tree.Arithmetic{Op: addOperator[op], Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) multiplicativeExpression() tree.Expression {
	left := ctx.unaryExpression()
	switch ctx.next() {
	case MUL, DIV, MOD:
		op, _ := ctx.consume()
		right := ctx.multiplicativeExpression()
		return tree.Arithmetic{Op: multOperator[op], Left: left, Right: right}
	}
	return left
}

func (ctx *parseContext) unaryExpression() tree.Expression {
	switch ctx.next() {
	case INV:
		ctx.consume()
		left := ctx.primary()
		return tree.BinaryNegation{left}
	case NOT:
		ctx.consume()
		left := ctx.primary()
		return tree.Negation{left}
	}
	return ctx.primary()
}

func (ctx *parseContext) collectArgs() []tree.Any {
	args := []tree.Any{}
	ctx.consume()
	for ctx.next() != RPAREN {
		args = append(args, ctx.logicalORExpression())
		switch ctx.next() {
		case RPAREN:
		case COMMA:
			ctx.consume()
		default:
			//TODO: error here
		}
	}
	ctx.consume()
	return args
}

func (ctx *parseContext) collectNumerics() []tree.Numeric {
	args := []tree.Numeric{}
	ctx.consume()
	for ctx.next() != RPAREN {
		args = append(args, ctx.logicalORExpression())
		switch ctx.next() {
		case RPAREN:
		case COMMA:
			ctx.consume()
		default:
			//TODO: error here
		}
	}
	ctx.consume()
	return args
}

func (ctx *parseContext) primary() tree.Expression {
	switch ctx.next() {
	case LPAREN:
		ctx.consume()
		val := ctx.logicalORExpression()
		op, _ := ctx.consume()
		if op != RPAREN {
			// TODO: raise error here
		}
		return val
	case ARG:
		_, data := ctx.consume()
		val, _ := strconv.Atoi(strings.TrimPrefix(string(data), "arg"))
		// This should never error out
		return tree.Argument{val}
	case IDENT:
		_, data := ctx.consume()
		if ctx.next() == LPAREN {
			return tree.Call{Name: string(data), Args: ctx.collectArgs()}
		}
		return tree.Variable{string(data)}
	case IN, NOTIN:
		op, _ := ctx.consume()
		if ctx.next() == LPAREN {
			all := ctx.collectNumerics()
			return tree.Inclusion{Positive: op == IN, Left: all[0], Rights: all[1:]}
		}
		// ERROR here
	case INT:
		_, data := ctx.consume()
		val, _ := strconv.ParseUint(string(data), 0, 32)
		return tree.NumericLiteral{uint32(val)}
	case TRUE:
		ctx.consume()
		return tree.BooleanLiteral{true}
	case FALSE:
		ctx.consume()
		return tree.BooleanLiteral{false}
	}

	// ERRROR here
	panic(fmt.Sprintf("Unexpected token: %s", tokens[ctx.next()]))
}
