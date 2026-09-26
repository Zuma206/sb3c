package lexer

import "fmt"

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
