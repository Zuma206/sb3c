package lexer

import (
	"errors"
	"fmt"
	"strings"
)

type Lexer struct {
	errors       []error
	tokens       []*Token
	errorSection *Section
	errorSrc     strings.Builder
	types        []*Type
	pos          Position
	src          string
	next         int
}

func NewLexer(src string, types []*Type) *Lexer {
	return &Lexer{
		errors:       make([]error, 0),
		tokens:       make([]*Token, 0),
		errorSection: nil,
		types:        types,
		pos: Position{
			LineNumber: 1,
			LineOffset: 1,
			Index:      0,
		},
		src:  src,
		next: 0,
	}
}

func (lexer *Lexer) newToken(tokenType *Type, len int) *Token {
	return &Token{
		Type: tokenType,
		Section: Section{
			Pos: lexer.pos,
			Src: lexer.src[lexer.pos.Index : lexer.pos.Index+len],
		},
	}
}

func (lexer *Lexer) getLongestMatch() (*Token, bool) {
	var token *Token
	for _, tokenType := range lexer.types {
		pos := tokenType.regex.FindStringSubmatchIndex(lexer.src[lexer.pos.Index:])
		if pos != nil && (token == nil || len(token.Src) < pos[1]) {
			token = lexer.newToken(tokenType, pos[1])
		}
	}
	if token == nil {
		return nil, false
	}
	return token, true
}

func (lexer *Lexer) consume(src string) {
	lexer.pos.LineOffset += len(src)
	lexer.pos.Index += len(src)
	for _, char := range src {
		if char == '\n' {
			lexer.pos.LineOffset = 1
			lexer.pos.LineNumber++
		}
	}
}

func (lexer *Lexer) consumeIntoError() {
	if lexer.errorSection == nil {
		lexer.errorSection = &Section{Pos: lexer.pos}
		lexer.errorSrc = strings.Builder{}
	}
	src := string(lexer.src[lexer.pos.Index])
	lexer.errorSrc.WriteString(src)
	lexer.consume(src)
}

var LexErr = errors.New("lex error")

func (lexer *Lexer) consumeError() {
	if lexer.errorSection != nil {
		lexer.errorSection.Src = lexer.errorSrc.String()
		err := fmt.Errorf("%w: %q %w", LexErr, lexer.errorSection.Src, &lexer.errorSection.Pos)
		lexer.errors = append(lexer.errors, err)
		lexer.errorSection = nil
	}
}

var EOFError = errors.New("eof")

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

func (lexer *Lexer) Next() (*Token, error) {
	token, err := lexer.Peek(0)
	if err != nil {
		return nil, err
	}
	lexer.next++
	return token, err
}

func (lexer *Lexer) parseAll() {
	var err error
	for err == nil {
		_, err = lexer.Next()
	}
}

func (lexer *Lexer) GetTokens() []*Token {
	lexer.parseAll()
	return lexer.tokens
}

func (lexer *Lexer) GetErrors() []error {
	lexer.parseAll()
	return lexer.errors
}
