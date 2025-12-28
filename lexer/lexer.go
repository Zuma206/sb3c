package lexer

import "errors"

// Lexes a file into a list of tokens and errors
type Lexer struct {
	// Lex errors encountered so far
	errors []*Section
	// Lex tokens parsed so far
	tokens []*Token
	// The types of tokens the lexer can parse
	types []*Type
	// The current working position of the lexer in the file
	pos Position
	// The source code being lexed
	src []byte
	// The index of the next token to be processed
	next int
}

func NewLexer(src []byte, types []*Type) *Lexer {
	return &Lexer{
		errors: make([]*Section, 0),
		tokens: make([]*Token, 0),
		types:  types,
		pos: Position{
			LineNumber: 1,
			LineOffset: 1,
			Index:      0,
		},
		src:  src,
		next: 0,
	}
}
