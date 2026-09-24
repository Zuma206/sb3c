package parser

// Represents a parse step
type Parse interface {
	Parse(*Parser) error
}

// Represents a parse step that constructs a value
type ParseValue[T any] interface {
	Parse
	ParseValue(*Parser) (T, error)
}

// A single-function version of the `Parse` interface
type ParseFunc func(*Parser) error

// Calls the underlying `ParseFunc`
func (parse ParseFunc) Parse(p *Parser) error {
	return parse(p)
}

// A single-function version of the `ParseValue` interface
type ParseValueFunc[T any] func(*Parser) (T, error)

// Calls the underlying `ParseValueFunc`
func (parseValue ParseValueFunc[T]) ParseValue(p *Parser) (T, error) {
	return parseValue(p)
}

// Calls the underlying `ParseValueFunc`, discarding the value
func (parseValue ParseValueFunc[T]) Parse(p *Parser) error {
	_, err := parseValue.ParseValue(p)
	return err
}
