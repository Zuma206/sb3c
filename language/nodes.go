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
	Members *utils.List[*Method]
}

type Method struct {
	Decorator *lexer.Token
	Name      *lexer.Token
	Args      *utils.List[*lexer.Token]
	Calls     *utils.List[*Call]
}

type Call struct {
	Path *lexer.Token
	Args *utils.List[*lexer.Token]
}
