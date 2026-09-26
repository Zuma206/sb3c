package parser

import "github.com/zuma206/sb3c/utils"

type ShouldConsumeCondition bool

var (
	ConsumeCondition       ShouldConsumeCondition = true
	SkipConsumingCondition ShouldConsumeCondition = false
)

func If[T any](ifCheck CheckableStepAny, consume ShouldConsumeCondition, stepIf Step[T], stepElse Step[T]) StepFunc[T] {
	var zeroValue T
	return func(p *Parser) (T, error) {
		if ifCheck.CanParse(p) == nil {
			if consume {
				if err := ifCheck.CanParse(p); err != nil {
					return zeroValue, err
				}
			}
			return stepIf.Parse(p)
		}
		return stepElse.Parse(p)
	}
}

func DoWhile[T any](doStep Step[T], whileCheck CheckableStepAny, consume ShouldConsumeCondition) StepFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for {
			value, err := doStep.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
			if whileCheck.CanParse(p) != nil {
				break
			}
			if consume {
				if _, err := whileCheck.ParseAny(p); err != nil {
					return nil, err
				}
			}
		}
		return list, nil
	}
}

func Until[T any](step Step[T], untilCheck CheckableStepAny, consume ShouldConsumeCondition) StepFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for untilCheck.CanParse(p) != nil {
			value, err := step.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		if consume {
			if _, err := untilCheck.ParseAny(p); err != nil {
				return nil, err
			}
		}
		return list, nil
	}
}

func While[T any](whileCheck CheckableStepAny, step Step[T], consume ShouldConsumeCondition) StepFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for whileCheck.CanParse(p) == nil {
			if consume {
				if _, err := whileCheck.ParseAny(p); err != nil {
					return nil, err
				}
			}
			value, err := step.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}

// Continually parses until the parser is finished
func UntilFinished[T any](step Step[T]) StepFunc[*utils.List[T]] {
	return func(p *Parser) (*utils.List[T], error) {
		list := utils.NewList[T]()
		for !p.Finished() {
			value, err := step.Parse(p)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		}
		return list, nil
	}
}
