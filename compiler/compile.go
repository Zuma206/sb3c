package compiler

import (
	"bytes"

	"github.com/zuma206/sb3c/codegen"
	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/sb3"
)

type CompileResult struct {
	Tokens  []*lexer.Token
	Program *language.Program
	Output  []byte
}

func Compile(src string) (*CompileResult, error) {
	l := lexer.NewLexer(src, language.Types)
	p := parser.NewParser(l.GetTokens())
	program, err := language.ParseProgram(p)
	if err != nil {
		return nil, err
	}
	sb3Project, err := codegen.Generate(program)
	if err != nil {
		return nil, err
	}
	return newResult(l.GetTokens(), program, sb3Project)
}

func newResult(tokens []*lexer.Token, program *language.Program, sb3Project *sb3.SB3) (*CompileResult, error) {
	var buf bytes.Buffer
	if _, err := sb3Project.WriteTo(&buf); err != nil {
		return nil, err
	}
	return &CompileResult{
		Output:  buf.Bytes(),
		Tokens:  tokens,
		Program: program,
	}, nil
}
