package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/utils"
)

type Parser[U any] struct {
	Debug  func(any)
	tokens []U
	index  int
}

func NewParser[U any](tokens []U) *Parser[U] {
	return &Parser[U]{
		tokens: tokens,
		index:  0,
	}
}

var (
	NegativeOffsetError = errors.New("negative offset")
	EOFError            = errors.New("eof")
)

func (parser *Parser[U]) Peek() (U, error) {
	if parser.index >= len(parser.tokens) {
		return utils.Zero[U](), fmt.Errorf("%w: no token left to peek", EOFError)
	}
	return parser.tokens[parser.index], nil
}

func (parser *Parser[U]) Consume() (U, error) {
	token, err := parser.Peek()
	if err != nil {
		return utils.Zero[U](), err
	}
	parser.index++
	if parser.Debug != nil {
		parser.Debug(token)
	}
	return token, nil
}

func (parser *Parser[U]) Finished() bool {
	return parser.index >= len(parser.tokens)
}
