package optional

import "errors"

type Optional[T any] struct {
	value T
	ok    bool
}

func Value[T any](v T) Optional[T] {
	return Optional[T]{
		value: v,
		ok:    true,
	}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

func Pointer[T any](p *T) Optional[T] {
	if p == nil {
		return Empty[T]()
	}

	return Value(*p)
}

func Omitzero[T comparable](v T) Optional[T] {
	var empty T

	if v == empty {
		return Empty[T]()
	}

	return Value(v)
}

func (o Optional[T]) Value() (T, bool) {
	return o.value, o.ok
}

func (o Optional[T]) IsEmpty() bool {
	return !o.ok
}

var ErrOptionalValueIsEmpty = errors.New("optional value is empty")

func (o Optional[T]) MustValue() T {
	if o.ok {
		return o.value
	}

	panic(ErrOptionalValueIsEmpty)
}
