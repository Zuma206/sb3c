package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Member struct {
	Decorators        *utils.List[*Call]
	Name              *lexer.Token
	AttributeOrMethod *AttributeOrMethod
}

var FailedMemberParseErr = errors.New("failed member parse")

var member = parser.Returns(func(member *Member) parser.StepAny {
	return parser.Err(FailedMemberParseErr,
		parser.Sequence(
			parser.Set(&member.Decorators,
				parser.While(parser.Token(Symbol, At),
					parser.Suffix(
						call,
						parser.Optional(parser.Type(Whitespace)),
					),
					parser.ConsumeCondition,
				),
			),
			parser.Optional(parser.Type(Whitespace)),
			parser.Set(&member.Name, parser.Type(Identifier)),
			parser.Optional(parser.Type(Whitespace)),
			parser.Set(&member.AttributeOrMethod, attributeOrMethod),
		),
	)
})

type AttributeOrMethod struct {
	Attribute *Attribute
	Method    *Method
}

var FailedAttributeOrMethodParseErr = errors.New("failed attribute or method parse")

var attributeOrMethod = parser.Returns(func(attributeOrMethod *AttributeOrMethod) parser.StepAny {
	return parser.Err(FailedAttributeOrMethodParseErr,
		parser.Switch(
			parser.Case(parser.Token(Symbol, OpenBracket),
				parser.Set(&attributeOrMethod.Method, method)),
			parser.Case(parser.OneOf(parser.Token(Symbol, Equals), parser.Token(Symbol, Semicolon)),
				parser.Set(&attributeOrMethod.Attribute, attribute))),
	)
})

type Attribute struct {
	Initializer *Expression
}

var attribute = parser.Returns(func(attribute *Attribute) parser.StepAny {
	return parser.Sequence(
		parser.If(parser.Token(Symbol, Equals), parser.ConsumeCondition,
			parser.Sequence(
				parser.Optional(parser.Type(Whitespace)),
				parser.Set(&attribute.Initializer, expression),
				parser.Optional(parser.Type(Whitespace)),
			),
			parser.Sequence(),
		),
		parser.Token(Symbol, Semicolon),
	)
})
