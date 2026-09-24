package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/lexer"
)

// Parses a single token by type, source, or both
type ParseToken struct {
	Type *lexer.Type
	Src  *string
}

func Type(tokenType *lexer.Type) *ParseToken {
	return &ParseToken{
		Type: tokenType,
	}
}

func Token(tokenType *lexer.Type, src string) *ParseToken {
	return &ParseToken{
		Type: tokenType,
		Src:  &src,
	}
}

var (
	InvalidTokenTypeErr = errors.New("invalid token type")
	InvalidTokenSrcErr  = errors.New("invalid token source")
)

// Parses a single token described by `ParseToken`, returning it
func (parseToken ParseToken) ParseValue(p *Parser) (*lexer.Token, error) {
	token, err := p.Peek(0)
	if err != nil {
		return nil, err
	}
	if parseToken.Type != nil && token.Type != parseToken.Type {
		err := fmt.Errorf("expected %q got %q %w", parseToken.Type.Name, token.Type.Name, &token.Pos)
		return nil, errors.Join(InvalidTokenTypeErr, err)
	}
	if parseToken.Src != nil && token.Src != *parseToken.Src {
		err := fmt.Errorf("expected %q got %q %w", *parseToken.Src, token.Src, &token.Pos)
		return nil, errors.Join(InvalidTokenSrcErr, err)
	}
	return p.Consume()
}

// Parses a single token described by `ParseToken` and then discards it
func (parseToken ParseToken) Parse(p *Parser) error {
	_, err := parseToken.ParseValue(p)
	return err
}
