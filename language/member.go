package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Member struct {
	Decorators        *utils.List[*Call]
	Name              *lexer.Token
	AttributeOrMethod *AttributeOrMethod
}

var member = parser.Value(func(member *Member) parser.Parse {
	return parser.All(
		parser.Store(&member.Decorators,
			parser.While(parser.Token(Symbol, At),
				parser.Affix(parser.All(),
					parser.Func(parseCall),
					parser.Optional(parser.Type(Whitespace))))),
		parser.Optional(parser.Type(Whitespace)),
		parser.Store(&member.Name, parser.Type(Identifier)),
		parser.Optional(parser.Type(Whitespace)),
		parser.Store(&member.AttributeOrMethod, attributeOrMethod),
	)
})

type AttributeOrMethod struct {
	Attribute *Attribute
	Method    *Method
}

var attributeOrMethod = parser.Value(func(attributeOrMethod *AttributeOrMethod) parser.Parse {
	return parser.Switch(
		parser.Case(parser.Token(Symbol, OpenBracket),
			parser.Store(&attributeOrMethod.Method, parser.Func(parseMethod))),
		parser.Case(parser.Token(Symbol, Equals),
			parser.Store(&attributeOrMethod.Attribute, attribute)))
})

type Attribute struct {
	Initializer *lexer.Token
}

var attribute = parser.Value(func(attribute *Attribute) parser.Parse {
	return parser.All(
		parser.Optional(parser.Type(Whitespace)),
		parser.Store(&attribute.Initializer, parser.Func(parseExpression)),
		parser.Optional(parser.Type(Whitespace)),
		parser.Token(Symbol, Semicolon),
	)
})
