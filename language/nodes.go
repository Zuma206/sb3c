package language

import "github.com/zuma206/sb3c/lexer"

type Program struct {
	Declarations []*ClassDeclaration
}

type ClassDeclaration struct {
	Name         *lexer.Token
	Super        *lexer.Token
	Declarations []any
}
