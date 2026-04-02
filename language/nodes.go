package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/utils"
)

type Program struct {
	Classes *utils.List[*Class]
}

type Class struct {
	Name    *lexer.Token
	Super   *lexer.Token
	Members *utils.List[*Member]
}

type Member struct {
	Decorators *utils.List[*Call]
	Name       *lexer.Token
	Value      MemberValue
}

type MemberValue struct {
	Attribute *Attribute
	Method    *Method
}

type Method struct {
	Args  *utils.List[*lexer.Token]
	Calls *utils.List[*Call]
}

type Attribute struct {
	Initializer *lexer.Token
}

type Call struct {
	Path *lexer.Token
	Args *utils.List[*lexer.Token]
}
