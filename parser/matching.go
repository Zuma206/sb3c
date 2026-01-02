package parser

import "github.com/zuma206/sb3c/lexer"

// A matcher takes in a token and performs validation/matching on it
type Matcher func(token *lexer.Token) bool

// Creates a matcher for a token's type
func Type(tokenType *lexer.Type) Matcher {
	return func(token *lexer.Token) bool {
		return token.Type == tokenType
	}
}
