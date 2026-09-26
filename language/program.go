package language

import (
	"github.com/zuma206/sb3c/lexer"
	. "github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Program struct {
	Classes *utils.List[*Class]
}

var ParseProgram = Returns(func(program *Program) StepAny[*lexer.Token] {
	return Set(&program.Classes,
		UntilFinished(Affix(Optional(Whitespace), class, Optional(Whitespace))),
	)
})
