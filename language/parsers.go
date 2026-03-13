package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
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
		return nil, err
	}
	classDeclaration.Declarations, err = parseMethods(p)
	if err != nil {
		return nil, err
	}
	return classDeclaration, nil
}

func parseMethods(p *parser.Parser) (*utils.List[*MethodDeclaration], error) {
	methods := utils.NewList[*MethodDeclaration]()
	for true {
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
		Args: utils.NewList[*lexer.Token](),
		Body: utils.NewList[*FunctionCall](),
	}
	if _, err := p.ConsumeIf(Symbol.WithSource(At)); err == nil {
		if method.Decorator, err = p.ConsumeIf(Identifier); err != nil {
			return nil, err
		}
	}
	if err := p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: Identifier, Result: &method.Name},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBracket)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(CloseBracket)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBrace)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(CloseBrace)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, err
	}
	return method, nil
}
