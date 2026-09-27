package language

import (
	"github.com/zuma206/sb3c/lexer"
)

var Types = lexer.MustBuildTypes(Symbol, Keyword, Identifier, Path, Whitespace, NumberLiteral, StringLiteral)

var (
	Symbol       = lexer.NewSet("Symbol")
	OpenBrace    = Symbol.New("{")
	CloseBrace   = Symbol.New("}")
	Equals       = Symbol.New("=")
	Semicolon    = Symbol.New(";")
	At           = Symbol.New("@")
	Period       = Symbol.New(".")
	OpenBracket  = Symbol.New("(")
	CloseBracket = Symbol.New(")")
	Comma        = Symbol.New(",")
)

var (
	Keyword      = lexer.NewSet("Keyword")
	ClassKeyword = Keyword.New("class")
	Extends      = Keyword.New("extends")
)

var (
	identifierRegex = `[A-Za-z$_][0-9a-zA-Z$_]*`
	pathRegex       = identifierRegex + `(\.` + identifierRegex + `)*`
)

var (
	NumberLiteral = lexer.NewRegex("NumberLiteral", `[0-9]([0-9_]*[0-9])?`)
	Identifier    = lexer.NewRegex("Identifier", identifierRegex)
	Path          = lexer.NewRegex("Path", pathRegex)
	Whitespace    = lexer.NewRegex("Whitespace", `[\n\r\t\ ]+`)
	StringLiteral = lexer.NewRegex("StringLiteral", `"([^"\\]|\\.)*"`)
)
