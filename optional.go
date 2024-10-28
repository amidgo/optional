package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

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

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.ok
}

var ErrOptionalValueIsEmpty = errors.New("optional value is empty")

func (o Optional[T]) MustGet() T {
	if o.ok {
		return o.value
	}

	panic(ErrOptionalValueIsEmpty)
}

func (o Optional[T]) IsEmpty() bool {
	return !o.ok
}

func (o Optional[T]) Pointer() *T {
	if o.ok {
		return &o.value
	}

	return nil
}

func (o Optional[T]) sqlNull() sql.Null[T] {
	return sql.Null[T]{
		V:     o.value,
		Valid: o.ok,
	}
}

func (o Optional[T]) Value() (driver.Value, error) {
	sqlNull := o.sqlNull()

	return sqlNull.Value()
}

func (o *Optional[T]) Scan(src any) error {
	sqlNull := sql.Null[T]{}

	err := sqlNull.Scan(src)
	if err != nil {
		return err
	}

	o.value = sqlNull.V
	o.ok = sqlNull.Valid

	return nil
}

var jsonNull = [4]byte{'n', 'u', 'l', 'l'}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}

	return json.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && [4]byte(data) == [4]byte{'n', 'u', 'l', 'l'} {
		*o = Empty[T]()

		return nil
	}

	var v T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = Value(v)

	return nil
}
