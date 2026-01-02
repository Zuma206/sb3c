package parser

import "github.com/zuma206/sb3c/lexer"

// Generic parser that parses through a slice of lex tokens
type Parser struct {
	tokens []*lexer.Token
	index  int
}

// Constructs a parser
func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		index:  0,
	}
}
