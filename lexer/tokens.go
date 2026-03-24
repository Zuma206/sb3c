package lexer

import "fmt"

// Marks a position in a file
type Position struct {
	// The byte index into the file
	Index int
	// The line number, starting at 1
	LineNumber int
	// The byte index into the line, starting at 1
	LineOffset int
}

// Implements the error interface, so positions can be wrapped inside of errors
func (position *Position) Error() string {
	return fmt.Sprintf("(@%d:%d)", position.LineNumber, position.LineOffset)
}

// Represents a section of source code in a file
type Section struct {
	// The start positon of the section
	Pos Position
	// The source code of the section
	Src string
}

// Represents a lex token
type Token struct {
	// The type of lex token it was parsed as
	Type *Type
	// The section of code this token is
	Section
}
