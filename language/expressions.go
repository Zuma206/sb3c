package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/parser"
)

func parseExpression(p *parser.Parser) (*lexer.Token, error) {
	return p.ConsumeIf(lexer.MatchAny(NumberLiteral, StringLiteral))
}
