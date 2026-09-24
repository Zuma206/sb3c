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

type AttributeOrMethod struct {
	Attribute *Attribute
	Method    *Method
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

var attributeOrMethod = parser.Value(func(attributeOrMethod *AttributeOrMethod) parser.Parse {
	return parser.Switch(
		parser.Case(parser.Token(Symbol, OpenBracket),
			parser.Store(&attributeOrMethod.Method, parser.Func(parseMethod))),
		parser.Case(parser.Token(Symbol, Equals),
			parser.Store(&attributeOrMethod.Attribute, parser.Func(parseAttribute))))
})

func parseAttribute(p *parser.Parser) (*Attribute, error) {
	attribute := &Attribute{}
	p.ConsumeIf(Whitespace)
	attribute.Initializer, _ = parseExpression(p)
	if err := p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(Semicolon)},
	}); err != nil {
		return nil, errors.Join(AttributeErr, err)
	}
	return attribute, nil
}
