package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
)

type Call struct {
	Path *lexer.Token
	Args []*Expression
}

var FailedCallParseErr = errors.New("failed call parse")

var call = Err(FailedCallParseErr,
	Returns(func(call *Call) StepAny[*lexer.Token] {
		return Sequence(
			Set(&call.Path, OneOf(Identifier, Path)),
			Optional(Whitespace), OpenBracket,
			Set(&call.Args,
				If(CloseBracket, SkipConsumingCondition,
					None[*Expression, *lexer.Token](), callArgs),
			),
			CloseBracket,
		)
	}),
)

var FailedCallArgsParseErr = errors.New("failed call args parse")

var callArgs = Err(FailedCallArgsParseErr,
	DoWhile(
		Affix(Optional(Whitespace), expression, Optional(Whitespace)),
		Comma, ConsumeCondition,
	),
)
