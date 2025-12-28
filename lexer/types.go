package lexer

import "regexp"

// Represents a type of token a lexer can encounter
type Type struct {
	// The type's human readable name
	Name  string
	regex *regexp.Regexp
}

// Creates a lex token type from it's human readable name, and a regex that matches it
func NewType(name string, regex string) *Type {
	return &Type{
		Name:  name,
		regex: regexp.MustCompile("^" + regex),
	}
}
