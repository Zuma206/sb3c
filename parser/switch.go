package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/utils"
)

// Represents a checkable case and a parse step to execute if that case is valid
type switchCase struct {
	canParse CanParseAny
	parse    ParseAny
}

// Creates a condition -> parse steps pairing that can be checked by a parser.Switch
func Case(canParse CanParseAny, parse ParseAny) *switchCase {
	return &switchCase{
		canParse: canParse,
		parse:    parse,
	}
}

var NoCaseHitErr = errors.New("no case hit")

// Checks the conditions of all cases and executes the first case to hit
func Switch(cases ...*switchCase) ParseFunc[utils.UnitType] {
	return func(p *Parser) (utils.UnitType, error) {
		for _, switchCase := range cases {
			if switchCase.canParse.CanParse(p) != nil {
				continue
			}
			_, err := switchCase.parse.ParseAny(p)
			return utils.Unit, err
		}
		token, err := p.Peek(0)
		if err != nil {
			return utils.Unit, err
		}
		err = fmt.Errorf("at token %s(%q, %w)", token.Type.Name, token.Src, &token.Pos)
		return utils.Unit, errors.Join(NoCaseHitErr, err)
	}
}
