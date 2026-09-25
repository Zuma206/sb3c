package language

import (
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

type Program struct {
	Classes *utils.List[*Class]
}

var ParseProgram = parser.Value(func(program *Program) parser.ParseAny {
	return parser.Store(&program.Classes,
		parser.UntilFinished(
			parser.Affix(
				parser.Optional(parser.Type(Whitespace)),
				class,
				parser.Optional(parser.Type(Whitespace)),
			)))
})
