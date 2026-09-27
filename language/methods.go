package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
)

type Method struct {
	Args  []*lexer.Token
	Calls []*Call
}

var FailedMethodParseErr = errors.New("failed method parse")

var method = Err(FailedMethodParseErr,
	Returns(func(method *Method) StepAny[*lexer.Token] {
		return Sequence(
			OpenBracket, Optional(Whitespace), CloseBracket,
			Optional(Whitespace),
			OpenBrace, Set(&method.Calls, methodCalls), CloseBrace,
		)
	}),
)

var FailedMethodCallsParseErr = errors.New("failed method calls parse")

var methodCalls = Err(FailedMethodCallsParseErr,
	Until(
		Affix(
			Optional(Whitespace), call,
			Sequence(Optional(Whitespace), Semicolon, Optional(Whitespace)),
		),
		CloseBrace,
		SkipConsumingCondition,
	),
)
