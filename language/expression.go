package language

import (
	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
)

type Expression struct {
	Token *lexer.Token
}

var expression = Returns(func(expression *Expression) StepAny[*lexer.Token] {
	return Set(&expression.Token, OneOf(NumberLiteral, StringLiteral))
})
