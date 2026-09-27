package parser

import (
	"errors"

	"github.com/zuma206/sb3c/utils"
)

func Set[T any, U any](result *T, step Step[T, U]) StepFunc[T, U] {
	return func(p *Parser[U]) (T, error) {
		value, err := step.Parse(p)
		if err != nil {
			return utils.Zero[T](), err
		}
		*result = value
		return value, err
	}
}

func Sequence[U any](steps ...StepAny[U]) StepFunc[utils.UnitType, U] {
	return func(p *Parser[U]) (utils.UnitType, error) {
		for _, parse := range steps {
			if _, err := parse.ParseAny(p); err != nil {
				return utils.Unit, err
			}
		}
		return utils.Unit, nil
	}
}

func Returns[T any, U any](f func(value *T) StepAny[U]) StepFunc[*T, U] {
	return func(p *Parser[U]) (*T, error) {
		var value T
		_, err := f(&value).ParseAny(p)
		if err == nil && p.Debug != nil {
			p.Debug(value)
		}
		return &value, err
	}
}

func Affix[T any, U any](prefix StepAny[U], step Step[T, U], suffix StepAny[U]) StepFunc[T, U] {
	return func(p *Parser[U]) (value T, err error) {
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

func Suffix[T any, U any](parse Step[T, U], suffix StepAny[U]) StepFunc[T, U] {
	return Affix(Sequence[U](), parse, suffix)
}

func Prefix[T any, U any](prefix StepAny[U], parse Step[T, U]) StepFunc[T, U] {
	return Affix(prefix, parse, Sequence[U]())
}

func Err[T any, U any](parentErr error, step Step[T, U]) StepFunc[T, U] {
	return func(p *Parser[U]) (T, error) {
		value, err := step.Parse(p)
		if err != nil {
			return value, errors.Join(parentErr, err)
		}
		return value, nil
	}
}

func None[T any, U any]() StepFunc[[]T, U] {
	return func(_ *Parser[U]) ([]T, error) {
		return []T{}, nil
	}
}

func Optional[T any, U any](step CheckableStep[T, U]) StepFunc[utils.UnitType, U] {
	return func(p *Parser[U]) (utils.UnitType, error) {
		if step.CanParse(p) == nil {
			_, err := step.Parse(p)
			return utils.Unit, err
		}
		return utils.Unit, nil
	}
}

type oneOf[T any, U any] []CheckableStep[T, U]

func OneOf[T any, U any](options ...CheckableStep[T, U]) oneOf[T, U] {
	return options
}

func (oneOf oneOf[T, U]) CanParse(p *Parser[U]) error {
	var err error
	for _, option := range oneOf {
		if err = option.CanParse(p); err == nil {
			return nil
		}
	}
	return err
}

func (oneOf oneOf[T, U]) Parse(p *Parser[U]) (T, error) {
	var err error
	for _, option := range oneOf {
		if err = option.CanParse(p); err == nil {
			return option.Parse(p)
		}
	}
	return utils.Zero[T](), err
}

func (oneOf oneOf[T, U]) ParseAny(p *Parser[U]) (any, error) {
	return oneOf.Parse(p)
}
