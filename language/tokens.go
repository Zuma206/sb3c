package language

import "github.com/zuma206/sb3c/lexer"

// String literal keywords & symbols
const (
	OpenBrace    = "{"
	CloseBrace   = "}"
	ClassKeyword = "class"
	Extends      = "extends"
	Equals       = "="
	Semicolon    = ";"
	At           = "@"
	Period       = "."
	OpenBracket  = "("
	CloseBracket = ")"
	Comma        = ","
)

// Token type sets
var (
	Symbols  = []string{OpenBrace, CloseBrace, Equals, Semicolon, At, Period, OpenBracket, CloseBracket, Comma}
	Keywords = []string{ClassKeyword, Extends}
)

// Token types
var (
	NumberLiteral = lexer.NewType("NumberLiteral", `[0-9]([0-9_]*[0-9])?`)
	Identifier    = lexer.NewType("Identifier", `[A-z$_][0-z$_]*`)
	Path          = lexer.NewType("Path", `[A-z$_](\.?[0-z$_])*`)
	Whitespace    = lexer.NewType("Whitespace", `[\n\r\t\ ]+`)
	Keyword       = lexer.NewTypeSet("Keyword", Keywords)
	Symbol        = lexer.NewTypeSet("Symbol", Symbols)
)

var Types = []*lexer.Type{Symbol, Keyword, Identifier, Path, Whitespace, NumberLiteral}
