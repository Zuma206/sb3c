package parser

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
