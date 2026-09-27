package lexer

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/zuma206/sb3c/parser"
)

type Type struct {
	Name  string
	regex *regexp.Regexp
}

var InvalidTokenTypeErr = errors.New("invalid token type")

func (tokenType *Type) CanParse(p *parser.Parser[*Token]) error {
	token, err := p.Peek()
	if err != nil {
		return err
	}
	if token.Type != tokenType {
		err := fmt.Errorf("expected %q, got %q %w", tokenType.Name, token.Type.Name, &token.Pos)
		return errors.Join(InvalidTokenTypeErr, err)
	}
	return nil
}

func (tokenType *Type) Parse(p *parser.Parser[*Token]) (*Token, error) {
	if err := tokenType.CanParse(p); err != nil {
		return nil, err
	}
	return p.Consume()
}

func (tokenType *Type) ParseAny(p *parser.Parser[*Token]) (any, error) {
	return tokenType.Parse(p)
}

type TypeBuilder interface {
	BuildType() (*Type, error)
}

func MustBuildTypes(typeBuilders ...TypeBuilder) []*Type {
	types := make([]*Type, len(typeBuilders))
	for i, typeBuilder := range typeBuilders {
		var err error
		types[i], err = typeBuilder.BuildType()
		if err != nil {
			panic(err)
		}
	}
	return types
}
