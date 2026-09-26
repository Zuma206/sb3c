package lexer

import (
	"regexp"
)

type Type struct {
	Name  string
	regex *regexp.Regexp
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
