package sb3

import "encoding/json"

type Input struct {
	shadow  string
	literal *Literal
}

func LiteralInput(literal *Literal) *Input {
	return &Input{literal: literal}
}

func ShadowInput(shadow string) *Input {
	return &Input{shadow: shadow}
}

type LiteralType uint8

var (
	LiteralNumber LiteralType = 4
	LiteralString LiteralType = 10
)

type Literal struct {
	Type  LiteralType
	Value string
}

func (literal *Literal) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{literal.Type, literal.Value})
}

func (input *Input) MarshalJSON() ([]byte, error) {
	if input.literal == nil {
		return json.Marshal([]any{1, input.shadow})
	}
	return json.Marshal([]any{1, input.literal})
}
