package parser

import "github.com/zuma206/sb3c/lexer"

// A matcher takes in a token and performs validation/matching on it
type Matcher func(token *lexer.Token) bool
