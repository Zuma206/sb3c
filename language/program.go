package language

import (
	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
)

type Program struct {
	Classes []*Class
}

var ParseProgram = Returns(func(program *Program) StepAny[*lexer.Token] {
	return Set(&program.Classes,
		UntilFinished(Affix(Optional(Whitespace), class, Optional(Whitespace))),
	)
})
