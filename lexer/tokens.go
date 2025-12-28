package lexer

// Marks a position in a file
type Position struct {
	// The byte index into the file
	Index int
	// The line number, starting at 1
	LineNumber int
	// The byte index into the line, starting at 1
	LineOffset int
}

// Represents a section of source code in a file
type Section struct {
	// The start positon of the section
	Pos Position
	// The source code of the section
	Src []byte
}

// Represents a lex token
type Token struct {
	// The type of lex token it was parsed as
	Type *Type
	// The section of code this token is
	Section
}
