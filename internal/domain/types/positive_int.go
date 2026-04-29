package domainTypes

import "fmt"

type Integer interface {
	int | int32 | int64 | int16 | int8
}

type PositiveInt[T Integer] struct {
	value T
}

func (n PositiveInt[T]) Value() T {
	return n.value
}

func NewPositiveInt[T Integer](value T) (PositiveInt[T], error) {
	if value <= 0 {
		return PositiveInt[T]{}, fmt.Errorf("value should be positive: %d", value)
	}

	return PositiveInt[T]{
		value: value,
	}, nil
}
