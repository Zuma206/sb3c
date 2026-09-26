package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
)

type Expression struct {
	Token *lexer.Token
}

var expression = parser.Returns(func(expression *Expression) parser.StepAny[*lexer.Token] {
	return parser.Set(
		&expression.Token,
		parser.OneOf(parser.Type(NumberLiteral), parser.Type(StringLiteral)),
	)
})
