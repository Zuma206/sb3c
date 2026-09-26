package lexer

import "regexp"

type Regex struct {
	expr string
	*Type
}

func NewRegex(name string, expr string) *Regex {
	return &Regex{
		expr: expr,
		Type: &Type{
			Name: name,
		},
	}
}

func (regex *Regex) BuildType() (*Type, error) {
	compiledRegex, err := regexp.Compile("^" + regex.expr)
	if err != nil {
		return nil, err
	}
	regex.Type.regex = compiledRegex
	return regex.Type, nil
}
