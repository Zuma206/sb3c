package lexer

import "github.com/zuma206/sb3c/utils"

// Lexes a file into a list of tokens and errors
type Lexer struct {
	// Lex errors encountered so far
	errors []*Section
	// Lex tokens parsed so far
	tokens []*Token
	// The current working position of the lexer in the file
	pos Position
	// The source code being lexed
	src []byte
}

func NewLexer(src []byte) *Lexer {
	return &Lexer{
		errors: make([]*Section, 0),
		tokens: make([]*Token, 0),
		pos: Position{
			LineNumber: 1,
			LineOffset: 1,
			Index:      0,
		},
		src: src,
	}
}
