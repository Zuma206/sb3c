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

var member = parser.Value(func(member *Member) parser.ParseAny {
	return parser.Err(FailedMemberParseErr,
		parser.All(
			parser.Store(&member.Decorators,
				parser.While(parser.Token(Symbol, At),
					parser.Suffix(
						call,
						parser.Optional(parser.Type(Whitespace)),
					),
					parser.ConsumeCondition,
				),
			),
			parser.Optional(parser.Type(Whitespace)),
			parser.Store(&member.Name, parser.Type(Identifier)),
			parser.Optional(parser.Type(Whitespace)),
			parser.Store(&member.AttributeOrMethod, attributeOrMethod),
		),
	)
})

type AttributeOrMethod struct {
	Attribute *Attribute
	Method    *Method
}

var FailedAttributeOrMethodParseErr = errors.New("failed attribute or method parse")

var attributeOrMethod = parser.Value(func(attributeOrMethod *AttributeOrMethod) parser.ParseAny {
	return parser.Err(FailedAttributeOrMethodParseErr,
		parser.Switch(
			parser.Case(parser.Token(Symbol, OpenBracket),
				parser.Store(&attributeOrMethod.Method, method)),
			parser.Case(parser.OneOf(parser.Token(Symbol, Equals), parser.Token(Symbol, Semicolon)),
				parser.Store(&attributeOrMethod.Attribute, attribute))),
	)
})

type Attribute struct {
	Initializer *Expression
}

var attribute = parser.Value(func(attribute *Attribute) parser.ParseAny {
	return parser.All(
		parser.If(parser.Token(Symbol, Equals), parser.ConsumeCondition,
			parser.All(
				parser.Optional(parser.Type(Whitespace)),
				parser.Store(&attribute.Initializer, expression),
				parser.Optional(parser.Type(Whitespace)),
			),
			parser.All(),
		),
		parser.Token(Symbol, Semicolon),
	)
})
