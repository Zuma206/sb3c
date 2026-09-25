package parser

import (
	"errors"
	"fmt"
)

// Represents a checkable case and a parse step to execute if that case is valid
type switchCase struct {
	canParse CanParse
	parse    Parse
}

// Creates a condition -> parse steps pairing that can be checked by a parser.Switch
func Case(canParse CanParse, parse Parse) *switchCase {
	return &switchCase{
		canParse: canParse,
		parse:    parse,
	}
}

var NoCaseHitErr = errors.New("no case hit")

// Checks the conditions of all cases and executes the first case to hit
func Switch(cases ...*switchCase) ParseFunc {
	return func(p *Parser) error {
		for _, switchCase := range cases {
			if switchCase.canParse.CanParse(p) != nil {
				continue
			}
			return switchCase.parse.Parse(p)
		}
		token, err := p.Peek(0)
		if err != nil {
			return err
		}
		err = fmt.Errorf("at token %s(%q, %w)", token.Type.Name, token.Src, &token.Pos)
		return errors.Join(NoCaseHitErr, err)
	}
}
