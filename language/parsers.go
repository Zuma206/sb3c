package language

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

var (
	ClassDeclarationError  = errors.New("class declaration")
	MethodDeclarationError = errors.New("method declaration")
)

// Parses a program (AST root)
func ParseProgram(p *parser.Parser) (*Program, error) {
	program := &Program{Classes: utils.NewList[*Class]()}
	for !p.Finished() {
		p.ConsumeIf(Whitespace)
		classDeclaration, err := parseClassDeclaration(p)
		if err != nil {
			return nil, err
		}
		program.Classes.PushBack(classDeclaration)
		p.ConsumeIf(Whitespace)
	}
	return program, nil
}

func parseClassDeclaration(p *parser.Parser) (*Class, error) {
	classDeclaration := &Class{Members: utils.NewList[*Method]()}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Keyword.WithSource(ClassKeyword)},
		{Matcher: Whitespace},
		{Matcher: Identifier, Result: &classDeclaration.Name},
		{Matcher: Whitespace},
		{Matcher: Keyword.WithSource(Extends)},
		{Matcher: Whitespace},
		{Matcher: Identifier, Result: &classDeclaration.Super},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBrace)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, fmt.Errorf("%w: %w", ClassDeclarationError, err)
	}
	classDeclaration.Members, err = parseMethods(p)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", MethodDeclarationError, err)
	}
	return classDeclaration, nil
}

func parseMethods(p *parser.Parser) (*utils.List[*Method], error) {
	methods := utils.NewList[*Method]()
	for true {
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(CloseBrace)); err == nil {
			break
		}
		method, err := parseMethod(p)
		if err != nil {
			return nil, err
		}
		methods.PushBack(method)
	}
	return methods, nil
}

func parseMethod(p *parser.Parser) (*Method, error) {
	method := &Method{
		Decorator: parseOptionalDecorator(p),
		Args:      utils.NewList[*lexer.Token](),
	}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: Identifier, Result: &method.Name},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBracket)},
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

func parseOptionalDecorator(p *parser.Parser) *lexer.Token {
	if _, err := p.ConsumeIf(Symbol.WithSource(At)); err == nil {
		if decorator, err := p.ConsumeIf(Identifier); err != nil {
			return decorator
		}
	}
	return nil
}

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
		functionCalls.PushBack(functionCall)
	}
	return functionCalls, nil
}

func parseCall(p *parser.Parser) (*Call, error) {
	call := &Call{}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Path, Result: &call.Path},
		{Matcher: Symbol.WithSource(OpenBracket)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, err
	}
	call.Args, err = parseCallArgs(p)
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Symbol.WithSource(CloseBracket)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(Semicolon)},
	}); err != nil {
		return nil, err
	}
	return call, nil
}

func parseCallArgs(p *parser.Parser) (*utils.List[*lexer.Token], error) {
	args := utils.NewList[*lexer.Token]()
	for true {
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
