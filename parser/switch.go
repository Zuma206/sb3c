package parser

import (
	"errors"

	"github.com/zuma206/sb3c/utils"
)

type switchCase[T any, U any] struct {
	check CheckableStepAny[U]
	step  Step[T, U]
}

func Case[T any, U any](check CheckableStepAny[U], step Step[T, U]) *switchCase[T, U] {
	return &switchCase[T, U]{
		check: check,
		step:  step,
	}
}

var NoCaseHitErr = errors.New("no case hit")

func Switch[T any, U any](cases ...*switchCase[T, U]) StepFunc[T, U] {
	return func(p *Parser[U]) (T, error) {
		for _, switchCase := range cases {
			if switchCase.check.CanParse(p) != nil {
				continue
			}
			return switchCase.step.Parse(p)
		}
		return utils.Zero[T](), NoCaseHitErr
	}
}
