package parser2

import (
	"strconv"

	"github.com/twtiger/gosecco/tree"
)

type parseContext struct {
	index  int
	tokens []tokenData
	atEnd  bool
}

// ARITH things to support:
// - Parenthesis
// - Binary and (&)
// - Binary or (|)
// - Binary xor (^)
// - Binary negation (~)
// - Left shift (<<)
// - Right shift (>>)
// - Modulo (%)

// BOOLEAN things to support:
// - Parenthesis
// - Boolean OR (||)
// - Boolean AND(&&)
// - Boolean negation (!)
// - Comparison operators
//   - Equal (==)
//   - Not equal (!=)
//   - Greater than (>)
//   - Greater or equal to (>=)
//   - Less than (<)
//   - Less than or equal to (<=)
//   - Bits set (this operator will mask the left hand side with the right hand, and return true if the result has any bits set) (&)
// - Inclusion:
//   in(arg0, 1,2,3,4)
//   notIn(arg0, 1, 2, 3, 4)

// OTHER things to support:
// - Variables
// - Calls

// TODO: none of these check errors or whatever

func parseExpression(expr string) tree.Expression {
	tokens, ok := tokenize(expr)
	if !ok {
		// TODO: errors etc
		return nil
	}
	ctx := parseContext{0, tokens, false}
	expression := ctx.expression()
	ctx.end()

	return expression
}

func (ctx *parseContext) next() token {
	if ctx.atEnd {
		return ILLEGAL
	}
	return ctx.tokens[ctx.index].t
}

func (ctx *parseContext) advance() {
	ctx.index++
	if ctx.index >= len(ctx.tokens) {
		ctx.atEnd = true
	}
}

func (ctx *parseContext) consume() (token, []byte) {
	if ctx.atEnd {
		return ILLEGAL, nil
	}
	res := ctx.tokens[ctx.index]
	ctx.advance()
	return res.t, res.td
}

func (ctx *parseContext) end() {
	if !ctx.atEnd {
		// Raise error here
	}
}

func (ctx *parseContext) expression() tree.Expression {
	term := ctx.term()
	if ctx.next() == ADD || ctx.next() == SUB {
		op, _ := ctx.consume()
		term2 := ctx.term()
		opx := tree.PLUS
		if op == SUB {
			opx = tree.MINUS
		}
		return tree.Arithmetic{Op: opx, Left: term, Right: term2}
	}
	return term
}

func (ctx *parseContext) term() tree.Expression {
	factor := ctx.factor()
	switch ctx.next() {
	case MUL, DIV:
		op, _ := ctx.consume()
		factor2 := ctx.factor()
		opx := tree.MULT
		if op == DIV {
			opx = tree.DIV
		}
		return tree.Arithmetic{Op: opx, Left: factor, Right: factor2}
	}
	return factor
}

func (ctx *parseContext) factor() tree.Expression {
	// TODO: here we can also have parenthesis and recursive stuff, arguments and other things
	_, data := ctx.consume()
	// TODO: check token type is actually integer here of course
	val, _ := strconv.ParseUint(string(data), 0, 32)
	return tree.NumericLiteral{uint32(val)}
}
