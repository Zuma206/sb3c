package language

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
	"github.com/zuma206/sb3c/utils"
)

func parseMethod(p *parser.Parser) (*Method, error) {
	method := &Method{Args: utils.NewList[*lexer.Token]()}
	var err error
	if err = p.Parse([]*parser.ParseStep{
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(CloseBracket)},
		{Matcher: Whitespace, Optional: true},
		{Matcher: Symbol.WithSource(OpenBrace)},
		{Matcher: Whitespace, Optional: true},
	}); err != nil {
		return nil, err
	}
	method.Calls, err = parseCalls(p)
	if err != nil {
		return nil, err
	}
	return method, nil
}

var CallSemicolonErr = errors.New("missing semicolon after call")

func parseCalls(p *parser.Parser) (*utils.List[*Call], error) {
	functionCalls := utils.NewList[*Call]()
	for true {
		p.ConsumeIf(Whitespace)
		if _, err := p.ConsumeIf(Symbol.WithSource(CloseBrace)); err == nil {
			break
		}
		functionCall, err := parseCall(p)
		if err != nil {
			return nil, err
		}
		if _, err := p.ConsumeIf(Symbol.WithSource(Semicolon)); err != nil {
			return nil, errors.Join(CallSemicolonErr, err)
		}
		functionCalls.PushBack(functionCall)
	}
	return functionCalls, nil
}
