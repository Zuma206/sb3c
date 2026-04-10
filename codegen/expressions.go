package codegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
)

var NonConstantExpressionErr = errors.New("non-constant expression")

func evaluateConstantExpression(token *lexer.Token) (any, error) {
	switch token.Type {
	case language.NumberLiteral:
		return strconv.Atoi(strings.ReplaceAll(token.Src, "_", ""))
	case language.StringLiteral:
		return parseStringLiteral(token.Src), nil
	default:
		err := fmt.Errorf("expected a non-constant expression, got %q %w", token.Src, &token.Pos)
		return nil, errors.Join(NonConstantExpressionErr, err)
	}
}

func parseStringLiteral(src string) string {
	return src[1 : len(src)-1]
}
