package compiler

import (
	"bytes"
	"errors"
	"io/fs"

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

func Compile(src string, fileSystem fs.FS) *CompileResult {
	result := &CompileResult{}
	l := lexer.NewLexer(src, language.Types)
	result.Tokens = l.GetTokens()
	if len(l.GetErrors()) != 0 {
		result.Err = errors.Join(l.GetErrors()...)
		return result
	}
	p := parser.NewParser(result.Tokens)
	result.Program, result.Err = language.ParseProgram(p)
	if result.Err != nil {
		return result
	}
	result.Output, result.Err = getOutput(result.Program, fileSystem)
	return result
}

func getOutput(program *language.Program, fileSystem fs.FS) ([]byte, error) {
	sb3Project, err := codegen.Generate(program, fileSystem)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err = sb3Project.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
