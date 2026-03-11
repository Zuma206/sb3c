package parser

import (
	"errors"
	"github.com/zuma206/sb3c/lexer"
)

// A matcher takes in a token and performs validation/matching on it
type MatcherFunc func(token *lexer.Token) error

func (matcher MatcherFunc) MatchLexToken(token *lexer.Token) error {
	return matcher(token)
}

type Matcher interface {
	MatchLexToken(*lexer.Token) error
}

var UnexpectedTokenTypeError = errors.New("unexpected token type")

var UnexpectedTokenSource = errors.New("unexpected token source")
