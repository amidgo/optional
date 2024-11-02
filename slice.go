package optional

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/lib/pq"
)

type SliceOptional[T any] struct {
	value []T
	ok    bool
}

func SliceEmpty[T any]() SliceOptional[T] {
	return SliceOptional[T]{}
}

func Slice[T any](value []T) SliceOptional[T] {
	return SliceOptional[T]{
		value: value,
		ok:    true,
	}
}

func SliceOmitZero[T any](value []T) SliceOptional[T] {
	return SliceOptional[T]{
		value: value,
		ok:    len(value) > 0,
	}
}

func SlicePointer[T any](value *[]T) SliceOptional[T] {
	if value == nil {
		return SliceEmpty[T]()
	}

	return Slice(*value)
}

func (s SliceOptional[T]) Get() ([]T, bool) {
	return s.value, s.ok
}

func (s SliceOptional[T]) GetCopy() ([]T, bool) {
	if !s.ok {
		return nil, false
	}

	dst := make([]T, 0, len(s.value))

	copy(dst, s.value)

	return dst, true
}

func (s SliceOptional[T]) MustGet() []T {
	if !s.ok {
		panic(ErrOptionalValueIsEmpty)
	}

	return s.value
}

func (s SliceOptional[T]) MustGetCopy() []T {
	if !s.ok {
		panic(ErrOptionalValueIsEmpty)
	}

	dst := make([]T, 0, len(s.value))

	copy(dst, s.value)

	return dst
}

func (s SliceOptional[T]) OmitZero() SliceOptional[T] {
	s.ok = len(s.value) > 0

	return s
}

func (s SliceOptional[T]) Pointer() *[]T {
	if !s.ok {
		return nil
	}

	return &s.value
}

func (s SliceOptional[T]) IsEmpty() bool {
	return !s.ok
}

func (o *SliceOptional[T]) Scan(src any) error {
	value := make([]T, 0)

	pqArray := pq.Array(&value)

	err := pqArray.Scan(value)
	if err != nil {
		return err
	}

	*o = Slice(value)

	return nil
}

func (o SliceOptional[T]) Value() (driver.Value, error) {
	if !o.ok {
		return nil, nil
	}

	return pq.Array(o.value).Value()
}

func (o SliceOptional[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return jsonNull[:], nil
	}

	return json.Marshal(o.value)
}

func (o *SliceOptional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && [4]byte(data) == jsonNull {
		*o = SliceEmpty[T]()

		return nil
	}

	var v []T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = Slice(v)

	return nil
}
