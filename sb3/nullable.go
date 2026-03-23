package sb3

import "encoding/json"

type Nullable[T any] struct {
	value   T
	nonNull bool
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{nonNull: false}
}

func NonNull[T any](value T) Nullable[T] {
	return Nullable[T]{value: value, nonNull: true}
}

func (nullable *Nullable[T]) MarshalJSON() ([]byte, error) {
	if nullable.nonNull {
		return json.Marshal(nullable.value)
	}
	return []byte("null"), nil
}
