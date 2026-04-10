package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

var MemberNameErr = errors.New("failed to parse member name")

func parseCommonMember(p *parser.Parser) (*Member, error) {
	member := &Member{}
	var err error
	member.Decorators, err = parseDecorators(p)
	if err != nil {
		return nil, err
	}
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: Identifier, Result: &member.Name},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, errors.Join(MemberNameErr, err)
	}
	return member, nil
}

var MemberSymbolErr = errors.New("invalid member symbol")

func parseMemberValue(p *parser.Parser) (MemberValue, error) {
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

func parseDecorators(p *parser.Parser) (*utils.List[*Call], error) {
	decorators := utils.NewList[*Call]()
	for {
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(At)); err != nil {
			break
		}
		decorator, err := parseCall(p)
		if err != nil {
			return nil, err
		}
		decorators.PushBack(decorator)
	}
	return decorators, nil
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
