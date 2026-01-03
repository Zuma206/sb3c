package compiler

import (
	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/visualisation"
)

func CompileFile(name string, src []byte) error {
	lex := lexer.NewLexer(src, language.Types)
	visualisation.Visualise(lex.GetErrors())
	visualisation.Visualise(lex.GetTokens())
	p := parser.NewParser(lex.GetTokens())
	program, err := language.ParseProgram(p)
	if err != nil {
		return err
	}
	visualisation.Visualise(program)
	return nil
}
