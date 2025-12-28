package compiler

import (
	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/visualisation"
)

func CompileFile(name string, src []byte) error {
	lex := lexer.NewLexer(src, language.Types)
	visualisation.Visualise(lex.GetTokens())
	return nil
}
