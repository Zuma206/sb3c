package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Method struct {
	Args  *utils.List[*lexer.Token]
	Calls *utils.List[*Call]
}

var method = parser.Value(func(method *Method) parser.ParseAny {
	return parser.All(
		parser.Token(Symbol, OpenBracket),
		parser.Store(&method.Args, args),
		parser.Token(Symbol, CloseBracket),
		parser.Optional(parser.Type(Whitespace)),
		parser.Token(Symbol, OpenBrace),
		parser.Store(&method.Calls, methodCalls),
		parser.Token(Symbol, CloseBrace),
	)
})

var args = parser.Until(
	parser.Affix(
		parser.Optional(parser.Type(Whitespace)),
		parser.Type(Identifier),
		parser.Optional(parser.Type(Whitespace)),
	),
	parser.Token(Symbol, CloseBracket),
)

var FailedMethodCallsParse = errors.New("failed method calls parse")

var methodCalls = parser.Until(
	parser.Suffix(call,
		parser.All(
			parser.Optional(parser.Type(Whitespace)),
			parser.Token(Symbol, Semicolon),
			parser.Optional(parser.Type(Whitespace)),
		),
	),
	parser.Token(Symbol, CloseBrace),
)
