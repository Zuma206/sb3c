package language

import "github.com/zuma206/sb3c/lexer"

type Program struct {
	Declarations []*ClassDeclaration
}

type ClassDeclaration struct {
	Name         *lexer.Token
	Super        *lexer.Token
	Declarations []*MethodDeclaration
}

type MethodDeclaration struct {
	Decorators []*lexer.Token
	Name       *lexer.Token
	Args       []*lexer.Token
	Body       []any
}
