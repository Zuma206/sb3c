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

var FailedMethodParseErr = errors.New("failed method parse")

var method = parser.Err(FailedMethodParseErr,
	parser.Returns(func(method *Method) parser.StepAny[*lexer.Token] {
		return parser.Sequence(
			parser.Token(Symbol, OpenBracket),
			parser.Optional(parser.Type(Whitespace)),
			parser.Token(Symbol, CloseBracket),
			parser.Optional(parser.Type(Whitespace)),
			parser.Token(Symbol, OpenBrace),
			parser.Set(&method.Calls, methodCalls),
			parser.Token(Symbol, CloseBrace),
		)
	}),
)

var FailedMethodCallsParseErr = errors.New("failed method calls parse")

var methodCalls = parser.Err(FailedMethodCallsParseErr,
	parser.Until(
		parser.Affix(
			parser.Optional(parser.Type(Whitespace)),
			call,
			parser.Sequence(
				parser.Optional(parser.Type(Whitespace)),
				parser.Token(Symbol, Semicolon),
				parser.Optional(parser.Type(Whitespace)),
			),
		),
		parser.Token(Symbol, CloseBrace),
		parser.SkipConsumingCondition,
	),
)
