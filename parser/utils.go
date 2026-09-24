package parser

import (
	"github.com/zuma206/sb3c/utils"
)

// Stores the value of a `ParseValue` step into a pointer
func Store[T any](result *T, parseValue ParseValue[T]) ParseFunc {
	return func(p *Parser) error {
		value, err := parseValue.ParseValue(p)
		if err != nil {
			return err
		}
		*result = value
		return nil
	}
}

// Parses all steps in sequence
func All(all ...Parse) ParseFunc {
	return func(p *Parser) error {
		for _, parse := range all {
			if err := parse.Parse(p); err != nil {
				return err
			}
		}
		return nil
	}
}

// Continually parses until the parser is finished
func UntilFinished[T any](parse ParseValue[T]) ParseValueFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for !p.Finished() {
			value, err := parse.ParseValue(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}

// Creates a value inline whilst parsing
func Value[T any](f func(value *T) Parse) ParseValueFunc[*T] {
	return func(p *Parser) (*T, error) {
		var value T
		return &value, f(&value).Parse(p)
	}
}

// Constructs a `ParseValueFunc[T]` with inferrance
func Func[T any](parseValueFunc ParseValueFunc[T]) ParseValueFunc[T] {
	return parseValueFunc
}

// Adds a prefix and suffix to a `ParseValue` whilst preserving the value
func Affix[T any](prefix Parse, parseValue ParseValue[T], suffix Parse) ParseValueFunc[T] {
	return func(p *Parser) (value T, err error) {
		if err = prefix.Parse(p); err != nil {
			return value, err
		}
		value, err = parseValue.ParseValue(p)
		if err != nil {
			return value, err
		}
		if err = suffix.Parse(p); err != nil {
			return value, err
		}
		return value, err
	}
}

// Parses a list of `T` until a specific token is reached and consumed
func Until[T any](parseValue ParseValue[T], parseToken *parseToken) ParseValueFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for {
			if err := parseToken.Parse(p); err != nil {
				break
			}
			value, err := parseValue.ParseValue(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}

// While the token described by `parseToken` can be consumed, `parseValue` will be parsed into a list
func While[T any](parseToken *parseToken, parseValue ParseValue[T]) ParseValueFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for parseToken.Parse(p) != nil {
			value, err := parseValue.ParseValue(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}
