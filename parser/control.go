package parser

import "github.com/zuma206/sb3c/utils"

// Determines if a conditional `canParse` should be checked, or checked and consumed
type ShouldConsumeCondition bool

var (
	// Consume the condition if it's true
	ConsumeCondition ShouldConsumeCondition = true
	// Check the condition but never consume it
	SkipConsumingCondition ShouldConsumeCondition = false
)

// Conditionally parses `parse` when `canParse` can be parsed
func If[T any](canParse CanParseAny, consume ShouldConsumeCondition, parseIf Parse[T], parseElse Parse[T]) ParseFunc[T] {
	var zeroValue T
	return func(p *Parser) (T, error) {
		if canParse.CanParse(p) == nil {
			if consume {
				if err := canParse.CanParse(p); err != nil {
					return zeroValue, err
				}
			}
			return parseIf.Parse(p)
		}
		return parseElse.Parse(p)
	}
}

// Parses `parse` into a list until, repeating until `canParse` can no longer be parsed
func DoWhile[T any](parse Parse[T], canParse CanParseAny, consume ShouldConsumeCondition) ParseFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for {
			value, err := parse.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
			if canParse.CanParse(p) != nil {
				break
			}
			if consume {
				if _, err := canParse.ParseAny(p); err != nil {
					return nil, err
				}
			}
		}
		return list, nil
	}
}

// Continuously parses `T` into a list until `canParse` can be parsed
func Until[T any](parse Parse[T], canParse CanParseAny, consume ShouldConsumeCondition) ParseFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for canParse.CanParse(p) != nil {
			value, err := parse.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		if consume {
			if _, err := canParse.ParseAny(p); err != nil {
				return nil, err
			}
		}
		return list, nil
	}
}

// Continually parses `parseValue` whilst `canParse` can be parsed
func While[T any](canParse CanParseAny, parse Parse[T], consume ShouldConsumeCondition) ParseFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for canParse.CanParse(p) == nil {
			if consume {
				if _, err := canParse.ParseAny(p); err != nil {
					return nil, err
				}
			}
			value, err := parse.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}

// Continually parses until the parser is finished
func UntilFinished[T any](parse Parse[T]) ParseFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for !p.Finished() {
			value, err := parse.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}
