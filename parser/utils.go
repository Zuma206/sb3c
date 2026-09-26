package parser

import (
	"errors"

	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/utils"
	"github.com/zuma206/sb3c/visualisation"
)

// Stores the value of a `Parse` step into a pointer
func Store[T any](result *T, parse Parse[T]) ParseFunc[T] {
	return func(p *Parser) (T, error) {
		value, err := parse.Parse(p)
		if err != nil {
			return value, err
		}
		*result = value
		return value, err
	}
}

// Parses all steps in sequence
func All(all ...ParseAny) ParseFunc[utils.UnitType] {
	return func(p *Parser) (utils.UnitType, error) {
		for _, parse := range all {
			if _, err := parse.ParseAny(p); err != nil {
				return utils.Unit, err
			}
		}
		return utils.Unit, nil
	}
}

// Creates a value inline whilst parsing
func Value[T any](f func(value *T) ParseAny) ParseFunc[*T] {
	return func(p *Parser) (*T, error) {
		var value T
		_, err := f(&value).ParseAny(p)
		return &value, err
	}
}

// Adds a prefix and suffix to a `ParseValue` whilst preserving the value
func Affix[T any](prefix ParseAny, parse Parse[T], suffix ParseAny) ParseFunc[T] {
	return func(p *Parser) (value T, err error) {
		if _, err = prefix.ParseAny(p); err != nil {
			return value, err
		}
		value, err = parse.Parse(p)
		if err != nil {
			return value, err
		}
		if _, err = suffix.ParseAny(p); err != nil {
			return value, err
		}
		return value, err
	}
}

// See `Affix`
func Suffix[T any](parse Parse[T], suffix ParseAny) ParseFunc[T] {
	return Affix(All(), parse, suffix)
}

// See `Affix`
func Prefix[T any](prefix ParseAny, parse Parse[T]) ParseFunc[T] {
	return Affix(prefix, parse, All())
}

// Joins any parse errors with a given parent error
func Err[T any](parentErr error, parse Parse[T]) ParseFunc[T] {
	return func(p *Parser) (T, error) {
		value, err := parse.Parse(p)
		if err != nil {
			return value, errors.Join(parentErr, err)
		}
		return value, nil
	}
}

// Logs a token as it's parsed
func Log(parse Parse[*lexer.Token]) ParseFunc[*lexer.Token] {
	return func(p *Parser) (*lexer.Token, error) {
		token, err := parse.Parse(p)
		if err != nil {
			return nil, err
		}
		visualisation.Visualise(token)
		return token, nil
	}
}
