package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

// Parses a program (AST root)
func ParseProgram(p *parser.Parser) (*Program, error) {
	program := &Program{Classes: utils.NewList[*Class]()}
	for !p.Finished() {
		p.ConsumeIf(Whitespace)
		class, err := parseClass(p)
		if err != nil {
			return nil, err
		}
		program.Classes.PushBack(class)
		p.ConsumeIf(Whitespace)
	}
	return program, nil
}

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

var AttributeErr = errors.New("failed to parse attribute")

func parseAttribute(p *parser.Parser) (*Attribute, error) {
	attribute := &Attribute{}
	if err := p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: NumberLiteral, Result: &attribute.Initializer, Optional: true},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(Semicolon)},
	}); err != nil {
		return nil, errors.Join(AttributeErr, err)
	}
	return attribute, nil
}

func parseMethod(p *parser.Parser) (*Method, error) {
	method := &Method{Args: utils.NewList[*lexer.Token]()}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(CloseBracket)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBrace)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, err
	}
	method.Calls, err = parseCalls(p)
	if err != nil {
		return nil, err
	}
	return method, nil
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

var CallSemicolonErr = errors.New("missing semicolon after call")

func parseCalls(p *parser.Parser) (*utils.List[*Call], error) {
	functionCalls := utils.NewList[*Call]()
	for true {
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(CloseBrace)); err == nil {
			break
		}
		functionCall, err := parseCall(p)
		if err != nil {
			return nil, err
		}
		if _, err := p.ConsumeIf(Symbol.WithSource(Semicolon)); err != nil {
			return nil, errors.Join(CallSemicolonErr, err)
		}
		functionCalls.PushBack(functionCall)
	}
	return functionCalls, nil
}

var (
	CallErr      = errors.New("failed to parse call")
	CallCloseErr = errors.New("failed to parse call close")
)

func parseCall(p *parser.Parser) (*Call, error) {
	call := &Call{}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Path, Result: &call.Path},
		{Matcher: Symbol.WithSource(OpenBracket)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, errors.Join(CallErr, err)
	}
	call.Args, err = parseCallArgs(p)
	if err != nil {
		return nil, err
	}
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Symbol.WithSource(CloseBracket)},
	}); err != nil {
		return nil, errors.Join(CallCloseErr, err)
	}
	return call, nil
}

func parseCallArgs(p *parser.Parser) (*utils.List[*lexer.Token], error) {
	args := utils.NewList[*lexer.Token]()
	for !p.Check(Symbol.WithSource(CloseBracket)) {
		p.ConsumeIf(Whitespace)
		arg, err := p.ConsumeIf(NumberLiteral)
		if err != nil {
			return nil, err
		}
		args.PushBack(arg)
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(Comma)); err != nil {
			break
		}
	}
	return args, nil
}
