package language

import (
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

// Parses a program (AST root)
func ParseProgram(p *parser.Parser) (*Program, error) {
	program := &Program{Classes: utils.NewList[*Class]()}
	for !p.Finished() {
		p.ConsumeIf(Whitespace)
		class, err := parseClass(p)
		if err != nil {
			return nil, err
		}
		program.Classes.PushBack(class)
		p.ConsumeIf(Whitespace)
	}
	return program, nil
}
