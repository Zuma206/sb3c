package sb3

import "encoding/json"

type Variable struct {
	Name  string
	Value any
}

func (variable *Variable) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{variable.Name, variable.Value})
}
