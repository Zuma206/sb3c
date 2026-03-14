package compiler

import (
	"os"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/sb3"
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
	outfile, err := os.OpenFile("project.sb3", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := sb3.NewSB3().WriteTo(outfile); err != nil {
		return err
	}
	return nil
}
