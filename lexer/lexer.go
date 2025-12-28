package lexer

import "errors"

// Lexes a file into a list of tokens and errors
type Lexer struct {
	// Lex errors encountered so far
	errors []*Section
	// Lex tokens parsed so far
	tokens []*Token
	// The current error being encountered by the lexer
	error *Section
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
		error:  nil,
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

// Constructs a new token with a given type, with the src of a given length
func (lexer *Lexer) newToken(tokenType *Type, len int) *Token {
	return &Token{
		Type: tokenType,
		Section: Section{
			Pos: lexer.pos,
			Src: lexer.src[lexer.pos.Index:len],
		},
	}
}

// Find the longest token match out of all the token types
func (lexer *Lexer) getLongestMatch() (*Token, bool) {
	var token *Token
	for _, tokenType := range lexer.types {
		pos := tokenType.regex.FindSubmatchIndex(lexer.src[lexer.pos.Index:])
		if pos != nil && (token.Src == nil || len(token.Src) < pos[1]) {
			token = lexer.newToken(tokenType, pos[1])
		}
	}
	if token == nil {
		return nil, false
	}
	return token, true
}

// Consumes source code, incrementing the lexer position past it
func (lexer *Lexer) consume(src []byte) {
	lexer.pos.LineOffset += len(src)
	lexer.pos.Index += len(src)
	for _, char := range src {
		if char == '\n' {
			lexer.pos.LineOffset = 1
			lexer.pos.LineNumber++
		}
	}
}

// Consumes source code into the current error
func (lexer *Lexer) consumeIntoError() {
	src := []byte{lexer.src[lexer.pos.Index]}
	if lexer.error == nil {
		lexer.error = &Section{Pos: lexer.pos, Src: src}
	} else {
		lexer.error.Src = append(lexer.error.Src, lexer.src[lexer.pos.Index])
	}
	lexer.consume(src)
}

// Consumes the current error into the errors slice
func (lexer *Lexer) consumeError() {
	if lexer.error != nil {
		lexer.errors = append(lexer.errors, lexer.error)
		lexer.error = nil
	}
}
