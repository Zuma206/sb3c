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
	program := &Program{Declarations: utils.NewList[*ClassDeclaration]()}
	for !p.Finished() {
		p.ConsumeIf(Whitespace)
		classDeclaration, err := parseClassDeclaration(p)
		if err != nil {
			return nil, err
		}
		program.Declarations.PushBack(classDeclaration)
		p.ConsumeIf(Whitespace)
	}
	return program, nil
}

func parseClassDeclaration(p *parser.Parser) (*ClassDeclaration, error) {
	classDeclaration := &ClassDeclaration{Declarations: utils.NewList[*MethodDeclaration]()}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Keyword.WithSource(Class)},
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
	classDeclaration.Declarations, err = parseMethods(p)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", MethodDeclarationError, err)
	}
	return classDeclaration, nil
}

func parseMethods(p *parser.Parser) (*utils.List[*MethodDeclaration], error) {
	methods := utils.NewList[*MethodDeclaration]()
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

func parseMethod(p *parser.Parser) (*MethodDeclaration, error) {
	method := &MethodDeclaration{
		Decorator: parseOptionalDecorator(p),
		Args:      utils.NewList[*lexer.Token](),
		Body:      utils.NewList[*FunctionCall](),
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
	method.Body, err = parseFunctionCalls(p)
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

func parseFunctionCalls(p *parser.Parser) (*utils.List[*FunctionCall], error) {
	functionCalls := utils.NewList[*FunctionCall]()
	for true {
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(CloseBrace)); err == nil {
			break
		}
		functionCall, err := parseFunctionCall(p)
		if err != nil {
			return nil, err
		}
		functionCalls.PushBack(functionCall)
	}
	return functionCalls, nil
}

func parseFunctionCall(p *parser.Parser) (*FunctionCall, error) {
	functionCall := &FunctionCall{}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Path, Result: &functionCall.Path},
		{Matcher: Symbol.WithSource(OpenBracket)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, err
	}
	functionCall.Args, err = parseCallArgs(p)
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Symbol.WithSource(CloseBracket)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(Semicolon)},
	}); err != nil {
		return nil, err
	}
	return functionCall, nil
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
