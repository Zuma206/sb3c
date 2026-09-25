package codegen

import (
	"errors"
	"strconv"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/sb3"
)

var (
	NonConstantExpressionErr                  = errors.New("non-constant expression")
	MissingInitializerOnConstantExpressionErr = errors.New("missing initializer on constant expression")
)

func expressionToConst(expression *language.Expression) (any, error) {
	if expression == nil {
		return nil, MissingInitializerOnConstantExpressionErr
	}
	switch {
	case expression.Token != nil:
		return tokenToConst(expression.Token)
	default:
		return nil, NonConstantExpressionErr
	}
}

func tokenToConst(token *lexer.Token) (any, error) {
	switch token.Type {
	case language.NumberLiteral:
		return strconv.ParseFloat(token.Src, 64)
	case language.StringLiteral:
		return parseStringLiteral(token.Src), nil
	default:
		panic("malformed single-token expression")
	}
}

func parseStringLiteral(src string) string {
	if len(src) < 2 {
		panic("malformed string token")
	}
	return src[1 : len(src)-1]
}

func expressionToInput(expression *language.Expression) *sb3.Input {
	switch {
	case expression.Token != nil:
		return tokenToInput(expression.Token)
	default:
		panic("malformed expression")
	}
}

func tokenToInput(token *lexer.Token) *sb3.Input {
	switch token.Type {
	case language.NumberLiteral:
		return sb3.LiteralInput(sb3.LiteralNumber, token.Src)
	case language.StringLiteral:
		return sb3.LiteralInput(sb3.LiteralString, parseStringLiteral(token.Src))
	default:
		panic("malformed single-token expression")
	}
}
