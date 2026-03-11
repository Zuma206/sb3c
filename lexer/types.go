package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Represents a type of token a lexer can encounter
type Type struct {
	// The type's human readable name
	Name  string
	regex *regexp.Regexp
}

var (
	UnexpectedTokenTypeError = errors.New("UnexpectedTokenTypeError")
)

// Makes Type implement Matcher interface, checking if the token matches the type
func (tokenType *Type) MatchLexToken(token *Token) error {
	if token.Type == tokenType {
		return nil
	}
	return fmt.Errorf("%w: expected %s, found %s",
		UnexpectedTokenTypeError, tokenType.Name, token.Type.Name)
}

// Creates a lex token type from it's human readable name, and a regex that matches it
func NewType(name string, regex string) *Type {
	return &Type{
		Name:  name,
		regex: regexp.MustCompile("^" + regex),
	}
}

// Creates a lex token type that matches to any one string from a set
func NewTypeSet(name string, set []string) *Type {
	var regex strings.Builder
	regex.WriteRune('(')
	for i, str := range set {
		if i > 0 {
			regex.WriteRune('|')
		}
		regex.WriteString(regexp.QuoteMeta(str))
	}
	regex.WriteRune(')')
	return NewType(name, regex.String())
}
