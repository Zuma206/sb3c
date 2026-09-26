package parser

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/utils"
	"github.com/zuma206/sb3c/visualisation"
)

func Set[T any](result *T, step Step[T]) StepFunc[T] {
	return func(p *Parser) (T, error) {
		value, err := step.Parse(p)
		if err != nil {
			return value, err
		}
		*result = value
		return value, err
	}
}

func Sequence(steps ...StepAny) StepFunc[utils.UnitType] {
	return func(p *Parser) (utils.UnitType, error) {
		for _, parse := range steps {
			if _, err := parse.ParseAny(p); err != nil {
				return utils.Unit, err
			}
		}
		return utils.Unit, nil
	}
}

func Returns[T any](f func(value *T) StepAny) StepFunc[*T] {
	return func(p *Parser) (*T, error) {
		var value T
		_, err := f(&value).ParseAny(p)
		return &value, err
	}
}

func Affix[T any](prefix StepAny, step Step[T], suffix StepAny) StepFunc[T] {
	return func(p *Parser) (value T, err error) {
		if _, err = prefix.ParseAny(p); err != nil {
			return value, err
		}
		value, err = step.Parse(p)
		if err != nil {
			return value, err
		}
		if _, err = suffix.ParseAny(p); err != nil {
			return value, err
		}
		return value, err
	}
}

func Suffix[T any](parse Step[T], suffix StepAny) StepFunc[T] {
	return Affix(Sequence(), parse, suffix)
}

func Prefix[T any](prefix StepAny, parse Step[T]) StepFunc[T] {
	return Affix(prefix, parse, Sequence())
}

func Err[T any](parentErr error, step Step[T]) StepFunc[T] {
	return func(p *Parser) (T, error) {
		value, err := step.Parse(p)
		if err != nil {
			return value, errors.Join(parentErr, err)
		}
		return value, nil
	}
}

func Log(step Step[*lexer.Token]) StepFunc[*lexer.Token] {
	return func(p *Parser) (*lexer.Token, error) {
		token, err := step.Parse(p)
		if err != nil {
			return nil, err
		}
		visualisation.Visualise(token)
		return token, nil
	}
}

func None[T any]() StepFunc[*utils.List[T]] {
	return func(_ *Parser) (*utils.List[T], error) {
		return utils.NewList[T](), nil
	}
}

func Optional[T any](step CheckableStep[T]) StepFunc[utils.UnitType] {
	return func(p *Parser) (utils.UnitType, error) {
		if step.CanParse(p) == nil {
			_, err := step.Parse(p)
			return utils.Unit, err
		}
		return utils.Unit, nil
	}
}
