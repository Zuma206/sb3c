package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Method struct {
	Args  *utils.List[*lexer.Token]
	Calls *utils.List[*Call]
}

var method = parser.Value(func(method *Method) parser.Parse {
	return parser.All(
		parser.Optional(parser.Type(Whitespace)),
		parser.Token(Symbol, OpenBracket),
		parser.Store(&method.Args, parser.Until(
			parser.Affix(
				parser.Optional(parser.Type(Whitespace)),
				parser.Type(Identifier),
				parser.Optional(parser.Type(Whitespace)),
			),
			parser.Token(Symbol, CloseBracket),
		)),
		parser.Token(Symbol, CloseBracket),
		parser.Optional(parser.Type(Whitespace)),
		parser.Token(Symbol, OpenBrace),
		parser.Store(&method.Calls,
			parser.Until(
				parser.Suffix(call,
					parser.All(
						parser.Optional(parser.Type(Whitespace)),
						parser.Token(Symbol, Semicolon),
						parser.Optional(parser.Type(Whitespace)),
					),
				),
				parser.Token(Symbol, CloseBrace),
			),
		),
		parser.Token(Symbol, CloseBrace),
	)
})
