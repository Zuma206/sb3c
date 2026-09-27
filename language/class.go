package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
)

type Class struct {
	Name    *lexer.Token
	Super   *lexer.Token
	Members []*Member
}

var FailedClassParseErr = errors.New("failed class parse")

var class = Returns(func(class *Class) StepAny[*lexer.Token] {
	return Err(FailedClassParseErr,
		Sequence(
			ClassKeyword, Whitespace, Set(&class.Name, Identifier),
			Whitespace, Extends, Whitespace, Set(&class.Super, Identifier),
			Optional(Whitespace), OpenBrace,
			Set(&class.Members, Until(
				Affix(Optional(Whitespace), member, Optional(Whitespace)),
				CloseBrace,
				ConsumeCondition,
			),
			),
		),
	)
})
