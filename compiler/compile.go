package compiler

import (
	"bytes"

	"github.com/zuma206/sb3c/codegen"
	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
)

type CompileResult struct {
	Tokens  []*lexer.Token
	Program *language.Program
	Output  []byte
	Err     error
}

func Compile(src string) *CompileResult {
	result := &CompileResult{}
	l := lexer.NewLexer(src, language.Types)
	result.Tokens = l.GetTokens()
	p := parser.NewParser(result.Tokens)
	result.Program, result.Err = language.ParseProgram(p)
	if result.Err != nil {
		return result
	}
	result.Output, result.Err = getOutput(result.Program)
	return result
}

func getOutput(program *language.Program) ([]byte, error) {
	sb3Project, err := codegen.Generate(program)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err = sb3Project.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
