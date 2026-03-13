package language

import (
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

// Parses a program (AST root)
func ParseProgram(p *parser.Parser) (*Program, error) {
	program := &Program{Declarations: utils.NewList[*ClassDeclaration]()}
	classDeclaration, err := parseClassDeclaration(p)
	if err != nil {
		return nil, err
	}
	program.Declarations.PushBack(classDeclaration)
	return program, nil
}

func parseClassDeclaration(p *parser.Parser) (*ClassDeclaration, error) {
	classDeclaration := &ClassDeclaration{Declarations: utils.NewList[*MethodDeclaration]()}
	if err := p.Parse([]*parser.ParseStep{
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
		{Matcher: Symbol.WithSource(CloseBrace)},
	}); err != nil {
		return nil, err
	}
	return classDeclaration, nil
}
