package language

import (
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/utils"
)

type Method struct {
	Args  *utils.List[*lexer.Token]
	Calls *utils.List[*Call]
}

type Call struct {
	Path *lexer.Token
	Args *utils.List[*lexer.Token]
}
