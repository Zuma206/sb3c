package lexer

import "github.com/zuma206/sb3c/utils"

type Lexer struct {
	errors *utils.List[*Section]
	tokens *utils.List[*Token]
	pos    Position
	src    []byte
}

func NewLexer(src []byte) *Lexer {
	return &Lexer{
		errors: utils.NewList[*Section](),
		tokens: utils.NewList[*Token](),
		pos: Position{
			LineNumber: 1,
			LineOffset: 1,
			Index:      0,
		},
		src: src,
	}
}
