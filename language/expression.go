package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
)

type Expression struct {
	Token *lexer.Token
}

var expression = parser.Value(func(expression *Expression) parser.ParseAny {
	return parser.Store(
		&expression.Token,
		parser.OneOf(parser.Type(NumberLiteral), parser.Type(StringLiteral)),
	)
})
