package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/utils"
)

type Program struct {
	Declarations *utils.List[*ClassDeclaration]
}

type ClassDeclaration struct {
	Name         *lexer.Token
	Super        *lexer.Token
	Declarations *utils.List[*MethodDeclaration]
}

type MethodDeclaration struct {
	Decorators *utils.List[*lexer.Token]
	Name       *lexer.Token
	Args       *utils.List[*lexer.Token]
	Body       *utils.List[*FunctionCall]
}

type FunctionCall struct {
	Name *lexer.Token
	Args *utils.List[*lexer.Token]
}
