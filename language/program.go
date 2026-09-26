package language

import (
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Program struct {
	Classes *utils.List[*Class]
}

var ParseProgram = parser.Returns(func(program *Program) parser.StepAny {
	return parser.Set(&program.Classes,
		parser.UntilFinished(
			parser.Affix(
				parser.Optional(parser.Type(Whitespace)),
				class,
				parser.Optional(parser.Type(Whitespace)),
			)))
})
