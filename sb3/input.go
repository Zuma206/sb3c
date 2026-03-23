package sb3

import "encoding/json"

type Input struct {
	shadow string
}

func (input *Input) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{1, input.shadow})
}
