package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

var (
	CallErr      = errors.New("failed to parse call")
	CallCloseErr = errors.New("failed to parse call close")
)

func parseCall(p *parser.Parser) (*Call, error) {
	call := &Call{}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: lexer.MatchAny(Path, Identifier), Result: &call.Path},
		{Matcher: Symbol.WithSource(OpenBracket)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, errors.Join(CallErr, err)
	}
	call.Args, err = parseCallArgs(p)
	if err != nil {
		return nil, err
	}
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Symbol.WithSource(CloseBracket)},
	}); err != nil {
		return nil, errors.Join(CallCloseErr, err)
	}
	return call, nil
}

func parseCallArgs(p *parser.Parser) (*utils.List[*lexer.Token], error) {
	args := utils.NewList[*lexer.Token]()
	for !p.Check(Symbol.WithSource(CloseBracket)) {
		p.ConsumeIf(Whitespace)
		arg, err := p.ConsumeIf(NumberLiteral)
		if err != nil {
			return nil, err
		}
		args.PushBack(arg)
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(Comma)); err != nil {
			break
		}
	}
	return args, nil
}
