package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/zuma206/sb3c/parser"
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

var InvalidTokenSourceErr = errors.New("invalid token source")

func (tokenParser *TokenParser) CanParse(p *parser.Parser[*Token]) error {
	if err := tokenParser.tokenType.CanParse(p); err != nil {
		return err
	}
	token, err := p.Peek()
	if err != nil {
		return err
	}
	if token.Src != tokenParser.src {
		err := fmt.Errorf("expected %q, got %q %w", tokenParser.src, token.Src, &token.Pos)
		return errors.Join(InvalidTokenSourceErr, err)
	}
	return nil
}

func (tokenParser *TokenParser) Parse(p *parser.Parser[*Token]) (*Token, error) {
	if err := tokenParser.CanParse(p); err != nil {
		return nil, err
	}
	return p.Consume()
}

func (tokenParser *TokenParser) ParseAny(p *parser.Parser[*Token]) (any, error) {
	return tokenParser.Parse(p)
}

func (set *Set) BuildType() (*Type, error) {
	var expr strings.Builder
	expr.WriteRune('(')
	for i, s := range set.strings {
		if i > 0 {
			expr.WriteRune('|')
		}
		expr.WriteString(regexp.QuoteMeta(s))
	}
	expr.WriteRune(')')
	regex, err := regexp.Compile("^" + expr.String())
	if err != nil {
		return nil, err
	}
	set.Type.regex = regex
	return set.Type, nil
}
