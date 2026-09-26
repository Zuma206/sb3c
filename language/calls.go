package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Call struct {
	Path *lexer.Token
	Args *utils.List[*Expression]
}

var FailedCallParseErr = errors.New("failed call parse")

var call = parser.Err(FailedCallParseErr,
	parser.Value(func(call *Call) parser.ParseAny {
		return parser.All(
			parser.Store(&call.Path,
				parser.OneOf(parser.Type(Identifier), parser.Type(Path)),
			),
			parser.Optional(parser.Type(Whitespace)),
			parser.Token(Symbol, OpenBracket),
			parser.Store(&call.Args, callArgs),
		)
	}),
)

var FailedCallArgsParseErr = errors.New("failed call args parse")

var callArgs = parser.Err(FailedCallArgsParseErr,
	parser.Until(
		parser.Affix(
			parser.Optional(parser.Type(Whitespace)),
			expression,
			parser.All(
				parser.Optional(parser.Type(Whitespace)),
				parser.Token(Symbol, Comma),
			),
		),
		parser.Token(Symbol, CloseBracket),
		parser.ConsumeCondition,
	),
)
