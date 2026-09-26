package parser

type Step[T any] interface {
	Parse(*Parser) (T, error)
	StepAny
}

type StepAny interface {
	ParseAny(*Parser) (any, error)
}

type CheckableStep[T any] interface {
	Step[T]
	CheckableStepAny
}

type CheckableStepAny interface {
	StepAny
	CanParse(*Parser) error
}

type StepFunc[T any] func(*Parser) (T, error)

func (parseFunc StepFunc[T]) Parse(p *Parser) (T, error) {
	return parseFunc(p)
}

func (parseFunc StepFunc[T]) ParseAny(p *Parser) (any, error) {
	return parseFunc.Parse(p)
}
