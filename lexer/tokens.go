package lexer

import (
	"fmt"
	"io"
)

type Position struct {
	Index      int
	LineNumber int
	LineOffset int
}

func (position *Position) Error() string {
	return fmt.Sprintf("%d:%d", position.LineNumber, position.LineOffset)
}

type Section struct {
	Pos Position
	Src string
}

type Token struct {
	Type *Type
	Section
}

func (token *Token) Visualise(w io.Writer) {
	fmt.Fprintf(w, "%s(%q, %d:%d)\n", token.Type.Name, token.Src, token.Pos.LineNumber, token.Pos.LineOffset)
}
