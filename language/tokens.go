package language

import "github.com/zuma206/sb3c/lexer"

// String literal keywords & symbols
const (
	OpenScope  = "{"
	CloseScope = "}"
	Class      = "class"
	Extends    = "extends"
)

// Token type sets
var (
	Whitespaces = []string{" ", "\n", "\r", "\t"}
	Symbols     = []string{OpenScope, CloseScope}
	Keywords    = []string{Class, Extends}
)

// Token types
var (
	Identifier = lexer.NewType("Identifier", `[A-z$_][0-z$_]*`)
	Whitespace = lexer.NewTypeSet("Whitespace", Whitespaces)
	Keyword    = lexer.NewTypeSet("Keyword", Keywords)
	Symbol     = lexer.NewTypeSet("Symbol", Symbols)
)

var Types = []*lexer.Type{Symbol, Keyword, Identifier, Whitespace}
