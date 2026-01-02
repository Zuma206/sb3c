package parser

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/lexer"
)

// A matcher takes in a token and performs validation/matching on it
type Matcher func(token *lexer.Token) error

var UnexpectedTokenTypeError = errors.New("unexpected token type")

// Creates a matcher for a token's type
func Type(tokenType *lexer.Type) Matcher {
	return func(token *lexer.Token) error {
		if token.Type != tokenType {
			return fmt.Errorf("%w: expected type %s, not %s",
				UnexpectedTokenTypeError, tokenType.Name, token.Type.Name)
		}
		return nil
	}
}

var UnexpectedTokenSource = errors.New("unexpected token source")

// Creates a matcher for a token's type and source code
func Token(tokenType *lexer.Type, source string) Matcher {
	matchType := Type(tokenType)
	src := []byte(source)
	return func(token *lexer.Token) error {
		if err := matchType(token); err != nil {
			return err
		}
		if !bytes.Equal(token.Src, src) {
			return fmt.Errorf("%w: expected source %q, not %q",
				UnexpectedTokenSource, src, token.Src)
		}
		return nil
	}
}
