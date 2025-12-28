package lexer

type Position struct {
	Index, LineNumber, LineOffset int
}

type Section struct {
	Pos Position
	Src []byte
}

type Token struct {
	Type *Type
	Section
}
