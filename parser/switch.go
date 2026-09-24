package parser

import (
	"errors"
	"fmt"
)

type switchCase struct {
	condition *parseToken
	parse     Parse
}

func Case(condition *parseToken, parse Parse) *switchCase {
	return &switchCase{
		condition: condition,
		parse:     parse,
	}
}

var NoCaseHitErr = errors.New("no case hit")

func Switch(cases ...*switchCase) ParseFunc {
	return func(p *Parser) error {
		for _, switchCase := range cases {
			if switchCase.condition.Parse(p) != nil {
				continue
			}
			return switchCase.parse.Parse(p)
		}
		token, err := p.Peek(0)
		if err != nil {
			return err
		}
		err = fmt.Errorf("token %s(%q, %w)", token.Type.Name, token.Src, &token.Pos)
		return errors.Join(NoCaseHitErr, err)
	}
}
