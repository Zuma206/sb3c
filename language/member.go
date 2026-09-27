package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Member struct {
	Decorators        *utils.List[*Call]
	Name              *lexer.Token
	AttributeOrMethod *AttributeOrMethod
}

var (
	FailedMemberParseErr = errors.New("failed member parse")
)

var member = Returns(func(member *Member) StepAny[*lexer.Token] {
	return Err(FailedMemberParseErr,
		Sequence(
			Set(&member.Decorators,
				decorators,
			),
			Optional(Whitespace),
			Set(&member.Name, Identifier),
			Optional(Whitespace),
			Set(&member.AttributeOrMethod, attributeOrMethod),
		),
	)
})

var FailedDecoratorsParseErr = errors.New("failed decorators parse")

var decorators = While(At,
	Suffix(call, Optional(Whitespace)),
	ConsumeCondition)

type AttributeOrMethod struct {
	Attribute *Attribute
	Method    *Method
}

var FailedAttributeOrMethodParseErr = errors.New("failed attribute or method parse")

var attributeOrMethod = Returns(func(attributeOrMethod *AttributeOrMethod) StepAny[*lexer.Token] {
	return Err(FailedAttributeOrMethodParseErr,
		Switch(
			Case(OpenBracket,
				Sequence(Set(&attributeOrMethod.Method, method))),
			Case(OneOf(Equals, Semicolon),
				Sequence(Set(&attributeOrMethod.Attribute, attribute)))),
	)
})

type Attribute struct {
	Initializer *Expression
}

var attribute = Returns(func(attribute *Attribute) StepAny[*lexer.Token] {
	return Sequence(
		If(Equals, ConsumeCondition,
			Sequence(
				Optional(Whitespace),
				Set(&attribute.Initializer, expression),
				Optional(Whitespace),
			),
			Sequence[*lexer.Token](),
		),
		Semicolon,
	)
})
