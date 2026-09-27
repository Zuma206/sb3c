package parser

type Step[T any, U any] interface {
	Parse(*Parser[U]) (T, error)
	StepAny[U]
}

type StepAny[U any] interface {
	ParseAny(*Parser[U]) (any, error)
}

type CheckableStep[T any, U any] interface {
	Step[T, U]
	CheckableStepAny[U]
}

type CheckableStepAny[U any] interface {
	StepAny[U]
	CanParse(*Parser[U]) error
}

type StepFunc[T any, U any] func(*Parser[U]) (T, error)

func (parseFunc StepFunc[T, U]) Parse(p *Parser[U]) (T, error) {
	return parseFunc(p)
}

func (parseFunc StepFunc[T, U]) ParseAny(p *Parser[U]) (any, error) {
	return parseFunc.Parse(p)
}
