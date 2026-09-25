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

func (parseToken *parseToken) getErr(token *lexer.Token) error {
	var (
		expectedType = "*"
		expectedSrc  = "*"
	)
	if parseToken.tokenType != nil {
		expectedType = parseToken.tokenType.Name
	}
	if parseToken.src != nil {
		expectedSrc = *parseToken.src
	}
	return fmt.Errorf("expected %s(%q) got %s(%q) %w",
		expectedType, expectedSrc, token.Type.Name, token.Src, &token.Pos)
}

// Checks if the token can be parsed from the current parser state
func (parseToken *parseToken) CanParse(p *Parser) error {
	token, err := p.Peek(0)
	if err != nil {
		return err
	}
	if parseToken.tokenType != nil && token.Type != parseToken.tokenType {
		return errors.Join(InvalidTokenTypeErr, parseToken.getErr(token))
	}
	if parseToken.src != nil && token.Src != *parseToken.src {
		return errors.Join(InvalidTokenSrcErr, parseToken.getErr(token))
	}
	return nil
}

// Parses a single token described by `ParseToken`, returning it
func (parseToken *parseToken) Parse(p *Parser) (*lexer.Token, error) {
	if err := parseToken.CanParse(p); err != nil {
		return nil, err
	}
	return p.Consume()
}

// Parses a single token described by `ParseToken` and then discards it
func (parseToken *parseToken) ParseAny(p *Parser) (any, error) {
	return parseToken.Parse(p)
}

type oneOf []*parseToken

func OneOf(parseToken *parseToken, parseTokens ...*parseToken) oneOf {
	return append(oneOf{parseToken}, parseTokens...)
}

func (oneOf oneOf) Parse(p *Parser) (*lexer.Token, error) {
	var err error
	for _, parseToken := range oneOf {
		var token *lexer.Token
		if token, err = parseToken.Parse(p); err == nil {
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

func (oneOf oneOf) ParseAny(p *Parser) (any, error) {
	return oneOf.Parse(p)
}
