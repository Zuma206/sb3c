package parser

import "github.com/zuma206/sb3c/utils"

type ShouldConsumeCondition bool

var (
	ConsumeCondition       ShouldConsumeCondition = true
	SkipConsumingCondition ShouldConsumeCondition = false
)

func If[T any, U any](
	ifCheck CheckableStepAny[U],
	consume ShouldConsumeCondition,
	stepIf Step[T, U], stepElse Step[T, U],
) StepFunc[T, U] {
	return func(p *Parser[U]) (T, error) {
		if ifCheck.CanParse(p) == nil {
			if consume {
				if err := ifCheck.CanParse(p); err != nil {
					return utils.Zero[T](), err
				}
			}
			return stepIf.Parse(p)
		}
		return stepElse.Parse(p)
	}
}

func DoWhile[T any, U any](
	doStep Step[T, U],
	whileCheck CheckableStepAny[U],
	consume ShouldConsumeCondition,
) StepFunc[[]T, U] {
	return func(p *Parser[U]) ([]T, error) {
		ts := []T{}
		for {
			value, err := doStep.Parse(p)
			if err != nil {
				return nil, err
			}
			ts = append(ts, value)
			if whileCheck.CanParse(p) != nil {
				break
			}
			if consume {
				if _, err := whileCheck.ParseAny(p); err != nil {
					return nil, err
				}
			}
		}
		return ts, nil
	}
}

func Until[T any, U any](
	step Step[T, U],
	untilCheck CheckableStepAny[U],
	consume ShouldConsumeCondition,
) StepFunc[[]T, U] {
	return func(p *Parser[U]) ([]T, error) {
		ts := []T{}
		for untilCheck.CanParse(p) != nil {
			value, err := step.Parse(p)
			if err != nil {
				return nil, err
			}
			ts = append(ts, value)
		}
		if consume {
			if _, err := untilCheck.ParseAny(p); err != nil {
				return nil, err
			}
		}
		return ts, nil
	}
}

func While[T any, U any](
	whileCheck CheckableStepAny[U],
	step Step[T, U],
	consume ShouldConsumeCondition,
) StepFunc[[]T, U] {
	return func(p *Parser[U]) ([]T, error) {
		ts := []T{}
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
			ts = append(ts, value)
		}
		return ts, nil
	}
}

func UntilFinished[T any, U any](step Step[T, U]) StepFunc[[]T, U] {
	return func(p *Parser[U]) ([]T, error) {
		ts := []T{}
		for !p.Finished() {
			value, err := step.Parse(p)
			if err != nil {
				return nil, err
			}
			ts = append(ts, value)
		}
		return ts, nil
	}
}
