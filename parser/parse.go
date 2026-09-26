package parser

import "github.com/zuma206/sb3c/utils"

// Represents a parse step that constructs a value
type Parse[T any] interface {
	Parse(*Parser) (T, error)
	ParseAny
}

type ParseAny interface {
	ParseAny(*Parser) (any, error)
}

// Represents a parse step that can be checked for errors before parsing.
// If an error would occur, it is returned, and the parser state is left unmodified.
type CanParse[T any] interface {
	Parse[T]
	CanParseAny
}

type CanParseAny interface {
	ParseAny
	CanParse(*Parser) error
}

// A single-function version of the `Parse` interface
type ParseFunc[T any] func(*Parser) (T, error)

// Calls the underlying `ParseFunc`
func (parseFunc ParseFunc[T]) Parse(p *Parser) (T, error) {
	return parseFunc(p)
}

// Calls `Parse` and casts to any
func (parseFunc ParseFunc[T]) ParseAny(p *Parser) (any, error) {
	return parseFunc.Parse(p)
}

// Parses an optional token
func Optional[T any](canParse CanParse[T]) ParseFunc[utils.UnitType] {
	return func(p *Parser) (utils.UnitType, error) {
		if canParse.CanParse(p) == nil {
			_, err := canParse.Parse(p)
			return utils.Unit, err
		}
		return utils.Unit, nil
	}
}
