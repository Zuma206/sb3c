package lexer

import (
	"errors"
	"fmt"
)

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
			Src: lexer.src[lexer.pos.Index : lexer.pos.Index+len],
		},
	}
}

// Find the longest token match out of all the token types
func (lexer *Lexer) getLongestMatch() (*Token, bool) {
	var token *Token
	for _, tokenType := range lexer.types {
		pos := tokenType.regex.FindSubmatchIndex(lexer.src[lexer.pos.Index:])
		if pos != nil && (token == nil || len(token.Src) < pos[1]) {
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

var EOFError = errors.New("eof")

// Parses the next token out of the source code, returning an error if the end of file is hit
func (lexer *Lexer) parseToken() (*Token, error) {
	for lexer.pos.Index < len(lexer.src) {
		token, ok := lexer.getLongestMatch()
		if ok {
			lexer.consume(token.Src)
			lexer.consumeError()
			return token, nil
		}
		lexer.consumeIntoError()
	}
	lexer.consumeError()
	return nil, EOFError
}

var IndexOutOfBounds = errors.New("index out of bounds")

// Returns the token at index i relative to the next token
func (lexer *Lexer) Peek(i int) (*Token, error) {
	index := lexer.next + i
	if index < 0 {
		return nil, fmt.Errorf("%w: token index %d is out of bounds", IndexOutOfBounds, i)
	}
	if index >= len(lexer.tokens) {
		for range index + 1 - len(lexer.tokens) {
			token, err := lexer.parseToken()
			if err != nil {
				return nil, err
			}
			lexer.tokens = append(lexer.tokens, token)
		}
	}
	return lexer.tokens[index], nil
}

// Returns the next token and increments the next value
func (lexer *Lexer) Next() (*Token, error) {
	token, err := lexer.Peek(0)
	if err != nil {
		return nil, err
	}
	lexer.next++
	return token, err
}

// Parses all of the file
func (lexer *Lexer) parseAll() {
	var err error
	for err == nil {
		_, err = lexer.Next()
	}
}

// Parse the remainder of the file and return all parsed tokens
func (lexer *Lexer) GetTokens() []*Token {
	lexer.parseAll()
	return lexer.tokens
}

// Parse the remainder of the file and return all errors
func (lexer *Lexer) GetErrors() []*Section {
	lexer.parseAll()
	return lexer.errors
}
