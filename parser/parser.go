package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/lexer"
)

type Parser struct {
	tokens []*lexer.Token
	index  int
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		index:  0,
	}
}

var (
	NegativeOffsetError = errors.New("negative offset")
	EOFError            = errors.New("eof")
)

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

func (parser *Parser) Consume() (*lexer.Token, error) {
	token, err := parser.Peek(0)
	if err != nil {
		return nil, err
	}
	parser.index++
	return token, nil
}

func (parser *Parser) Finished() bool {
	return parser.index >= len(parser.tokens)
}
