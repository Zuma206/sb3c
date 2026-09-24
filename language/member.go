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
	AttributeOrMethod AttributeOrMethod
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
		parser.Store(&member.AttributeOrMethod, parser.Func(parseAttributeOrMethod)),
	)
})

var MemberSymbolErr = errors.New("invalid member symbol")

func parseAttributeOrMethod(p *parser.Parser) (AttributeOrMethod, error) {
	symbol, err := p.ConsumeIf(lexer.MatchAny(
		Symbol.WithSource(Equals), Symbol.WithSource(Semicolon), Symbol.WithSource(OpenBracket)))
	value := MemberValue{}
	if err != nil {
		return value, errors.Join(MemberSymbolErr, err)
	}
	switch symbol.Src {
	case Equals:
		value.Attribute, err = parseAttribute(p)
	case Semicolon:
		value.Attribute = &Attribute{}
	case OpenBracket:
		value.Method, err = parseMethod(p)
	default:
		// If this panic triggers, check the switch statement has a case for every lexer.MatchAny param
		panic("class member parsed invalid symbol as correct")
	}
	return value, err
}

var AttributeErr = errors.New("failed to parse attribute")

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
