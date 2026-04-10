package language

import (
	"errors"

	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

var ClassHeaderError = errors.New("class header error")

func parseClass(p *parser.Parser) (*Class, error) {
	class := &Class{}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Keyword.WithSource(ClassKeyword)},
		{Matcher: Whitespace},
		{Matcher: Identifier, Result: &class.Name},
		{Matcher: Whitespace},
		{Matcher: Keyword.WithSource(Extends)},
		{Matcher: Whitespace},
		{Matcher: Identifier, Result: &class.Super},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBrace)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, errors.Join(ClassHeaderError, err)
	}
	class.Members, err = parseMembers(p)
	if err != nil {
		return nil, err
	}
	return class, nil
}

func parseMembers(p *parser.Parser) (*utils.List[*Member], error) {
	members := utils.NewList[*Member]()
	for {
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(CloseBrace)); err == nil {
			break
		}
		member, err := parseCommonMember(p)
		if err != nil {
			return nil, err
		}
		member.Value, err = parseMemberValue(p)
		if err != nil {
			return nil, err
		}
		members.PushBack(member)
	}
	return members, nil
}
