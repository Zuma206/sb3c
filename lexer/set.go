package lexer

import (
	"regexp"
	"strings"
)

type Set struct {
	strings []string
	*Type
}

func NewSet(name string) *Set {
	return &Set{
		strings: []string{},
		Type: &Type{
			Name: name,
		},
	}
}

type TokenParser struct {
	tokenType *Type
	src       string
}

func (set *Set) New(src string) *TokenParser {
	set.strings = append(set.strings, src)
	return &TokenParser{
		tokenType: set.Type,
		src:       src,
	}
}

func (set *Set) BuildType() (*Type, error) {
	var expr strings.Builder
	expr.WriteRune('(')
	for i, s := range set.strings {
		if i != 0 {
			expr.WriteRune('|')
		}
		expr.WriteString(regexp.QuoteMeta(s))
	}
	expr.WriteRune(')')
	regex, err := regexp.Compile(expr.String())
	if err != nil {
		return nil, err
	}
	set.Type.regex = regex
	return set.Type, nil
}
