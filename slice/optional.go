package sliceoptional

import (
	"encoding/json"

	"github.com/amidgo/optional"
)

type Slice[T any] struct {
	value []T
	ok    bool
}

func New[T any](value []T) Slice[T] {
	return Slice[T]{
		value: value,
		ok:    true,
	}
}

func Pointer[T any](value *[]T) Slice[T] {
	if value == nil {
		return Slice[T]{}
	}

	return New(*value)
}

func (s Slice[T]) Get() ([]T, bool) {
	return s.value, s.ok
}

func (s Slice[T]) GetCopy() ([]T, bool) {
	if !s.ok {
		return nil, false
	}

	dst := make([]T, 0, len(s.value))

	copy(dst, s.value)

	return dst, true
}

func (s Slice[T]) MustGet() []T {
	v, ok := s.Get()
	if !ok {
		panic(optional.ErrOptionalValueIsEmpty)
	}

	return v
}

func (s Slice[T]) MustGetCopy() []T {
	v, ok := s.GetCopy()
	if !ok {
		panic(optional.ErrOptionalValueIsEmpty)
	}

	return v
}

func (s Slice[T]) OmitZero() Slice[T] {
	s.ok = len(s.value) > 0

	return s
}

func (s Slice[T]) Pointer() *[]T {
	if !s.ok {
		return nil
	}

	return &s.value
}

func (s Slice[T]) IsEmpty() bool {
	return !s.ok
}

func (s Slice[T]) IsZero() bool {
	return !s.ok
}

var jsonNull = [4]byte{'n', 'u', 'l', 'l'}

func (o Slice[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return jsonNull[:], nil
	}

	return json.Marshal(o.value)
}

func (o *Slice[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && [4]byte(data) == jsonNull {
		*o = Slice[T]{}

		return nil
	}

	var v []T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = New(v)

	return nil
}
