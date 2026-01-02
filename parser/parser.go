package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/lexer"
)

// Generic parser that parses through a slice of lex tokens
type Parser struct {
	tokens []*lexer.Token
	index  int
}

// Constructs a parser
func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		index:  0,
	}
}

var (
	// Occurs when a negative offset is passed into Peek
	NegativeOffsetError = errors.New("negative offset")
	// Occurs when a Peek is called that peeks past the end of the token list
	EOFError = errors.New("eof")
)

// Returns the token at the given offset from the parse index
func (parser *Parser) Peek(offset int) (*lexer.Token, error) {
	if offset < 0 {
		return nil, fmt.Errorf("%w: %d is less than zero", NegativeOffsetError, offset)
	}
	index := parser.index + offset
	if index >= len(parser.tokens) {
		return nil, fmt.Errorf("%w: index %d is out of bounds", EOFError, index)
	}
	return parser.tokens[parser.index], nil
}
