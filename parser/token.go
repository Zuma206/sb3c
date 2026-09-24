package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/lexer"
)

// Parses a single token by type, source, or both
type parseToken struct {
	tokenType *lexer.Type
	src       *string
}

func Type(tokenType *lexer.Type) *parseToken {
	return &parseToken{
		tokenType: tokenType,
	}
}

func Token(tokenType *lexer.Type, src string) *parseToken {
	return &parseToken{
		tokenType: tokenType,
		src:       &src,
	}
}

var (
	InvalidTokenTypeErr = errors.New("invalid token type")
	InvalidTokenSrcErr  = errors.New("invalid token source")
)

// Checks if the token can be parsed from the current parser state
func (parseToken parseToken) CanParse(p *Parser) error {
	token, err := p.Peek(0)
	if err != nil {
		return err
	}
	if parseToken.tokenType != nil && token.Type != parseToken.tokenType {
		err := fmt.Errorf("expected %q got %q %w", parseToken.tokenType.Name, token.Type.Name, &token.Pos)
		return errors.Join(InvalidTokenTypeErr, err)
	}
	if parseToken.src != nil && token.Src != *parseToken.src {
		err := fmt.Errorf("expected %q got %q %w", *parseToken.src, token.Src, &token.Pos)
		return errors.Join(InvalidTokenSrcErr, err)
	}
	return nil
}

// Parses a single token described by `ParseToken`, returning it
func (parseToken parseToken) ParseValue(p *Parser) (*lexer.Token, error) {
	if err := parseToken.CanParse(p); err != nil {
		return nil, err
	}
	return p.Consume()
}

// Parses a single token described by `ParseToken` and then discards it
func (parseToken parseToken) Parse(p *Parser) error {
	_, err := parseToken.ParseValue(p)
	return err
}

type oneOf []*parseToken

func OneOf(parseToken *parseToken, parseTokens ...*parseToken) oneOf {
	return append(oneOf{parseToken}, parseTokens...)
}

func (oneOf oneOf) ParseValue(p *Parser) (*lexer.Token, error) {
	var err error
	for _, parseToken := range oneOf {
		var token *lexer.Token
		if token, err = parseToken.ParseValue(p); err == nil {
			return token, err
		}
	}
	return nil, err
}

func (oneOf oneOf) CanParse(p *Parser) error {
	var err error
	for _, parseToken := range oneOf {
		if err = parseToken.CanParse(p); err == nil {
			return nil
		}
	}
	return err
}

func (oneOf oneOf) Parse(p *Parser) error {
	_, err := oneOf.ParseValue(p)
	return err
}
