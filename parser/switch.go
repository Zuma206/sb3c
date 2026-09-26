package parser

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/utils"
)

type switchCase struct {
	canParse CheckableStepAny
	parse    StepAny
}

func Case(check CheckableStepAny, step StepAny) *switchCase {
	return &switchCase{
		canParse: check,
		parse:    step,
	}
}

var NoCaseHitErr = errors.New("no case hit")

func Switch(cases ...*switchCase) StepFunc[utils.UnitType] {
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
