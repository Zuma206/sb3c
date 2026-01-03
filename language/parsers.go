package language

import (
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

// Parses a program (AST root)
func ParseProgram(p *parser.Parser) (*Program, error) {
	program := &Program{Declarations: utils.NewList[*ClassDeclaration]()}
	return program, nil
}
