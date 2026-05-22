package valueobject

import "fmt"

type LiteralOnly interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | complex64 | complex128 | string
}

type LiteralBase[T LiteralOnly] struct {
	v T
}

func (l LiteralBase[T]) Value() T {
	return l.v
}

func (l LiteralBase[T]) String() string {
	return fmt.Sprintf("%v", l.v)
}

func (l LiteralBase[T]) Equals(o LiteralBase[T]) bool {
	return l.v == o.v
}
