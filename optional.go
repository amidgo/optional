package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Optional[T comparable] struct {
	value T
	ok    bool
}

func Empty[T comparable]() Optional[T] {
	return Optional[T]{}
}

func Comparable[T comparable](v T) Optional[T] {
	return Optional[T]{
		value: v,
		ok:    true,
	}
}

func OmitZero[T comparable](v T) Optional[T] {
	var zeroValue T

	if v == zeroValue {
		return Empty[T]()
	}

	return Comparable(v)
}

func Pointer[T comparable](p *T) Optional[T] {
	if p == nil {
		return Empty[T]()
	}

	return Comparable(*p)
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.ok
}

func (o Optional[T]) MustGet() T {
	if !o.ok {
		panic(ErrOptionalValueIsEmpty)
	}

	return o.value
}

func (o Optional[T]) OmitZero() Optional[T] {
	var zeroValue T

	o.ok = o.value != zeroValue

	return o
}

func (o Optional[T]) Pointer() *T {
	if o.ok {
		return &o.value
	}

	return nil
}

func (o Optional[T]) IsEmpty() bool {
	return !o.ok
}

func (o Optional[T]) IsZero() bool {
	return !o.ok
}

func (o Optional[T]) sqlNull() sql.Null[T] {
	return sql.Null[T]{
		V:     o.value,
		Valid: o.ok,
	}
}

func (o Optional[T]) Value() (driver.Value, error) {
	return o.sqlNull().Value()
}

func (o *Optional[T]) Scan(src any) error {
	sqlNull := sql.Null[T]{}

	err := sqlNull.Scan(src)
	if err != nil {
		return err
	}

	*o = Optional[T]{
		value: sqlNull.V,
		ok:    sqlNull.Valid,
	}

	return nil
}

var jsonNull = [4]byte{'n', 'u', 'l', 'l'}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return jsonNull[:], nil
	}

	return json.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && [4]byte(data) == jsonNull {
		*o = Empty[T]()

		return nil
	}

	var v T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = Comparable(v)

	return nil
}
