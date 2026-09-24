package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Class struct {
	Name    *lexer.Token
	Super   *lexer.Token
	Members *utils.List[*Member]
}

var FailedClassParseErr = errors.New("failed class parse")

var class = parser.Value(func(class *Class) parser.Parse {
	return parser.Err(FailedClassParseErr,
		parser.All(
			parser.Token(Keyword, ClassKeyword),
			parser.Type(Whitespace),
			parser.Store(&class.Name, parser.Type(Identifier)),
			parser.Type(Whitespace),
			parser.Token(Keyword, Extends),
			parser.Type(Whitespace),
			parser.Store(&class.Super, parser.Type(Identifier)),
			parser.Optional(parser.Type(Whitespace)),
			parser.Token(Symbol, OpenBrace),
			parser.Store(&class.Members,
				parser.Until(
					parser.Affix(
						parser.Optional(parser.Type(Whitespace)),
						member,
						parser.Optional(parser.Type(Whitespace)),
					),
					parser.Token(Symbol, CloseBrace),
				),
			),
		),
	)
})
