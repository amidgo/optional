package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Optional[T any] struct {
	value T
	ok    bool
}

func New[T any](v T) Optional[T] {
	return Optional[T]{
		value: v,
		ok:    true,
	}
}

func Pointer[T any](p *T) Optional[T] {
	if p == nil {
		return Optional[T]{}
	}

	return New(*p)
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.ok
}

func (o Optional[T]) MustGet() T {
	v, ok := o.Get()
	if !ok {
		panic(ErrOptionalValueIsEmpty)
	}

	return v
}

func (o Optional[T]) Pointer() *T {
	if o.ok {
		return &o.value
	}

	return nil
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
	if o.IsZero() {
		return jsonNull[:], nil
	}

	return json.Marshal(o.value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && [4]byte(data) == jsonNull {
		*o = Optional[T]{}

		return nil
	}

	var v T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = New(v)

	return nil
}

func OmitZero[T comparable](op Optional[T]) Optional[T] {
	var zeroValue T

	op.ok = op.value != zeroValue

	return op
}

func Convert[T, O any](op Optional[T], f func(T) O) Optional[O] {
	v, ok := op.Get()
	if !ok {
		return Optional[O]{}
	}

	return New(f(v))
}
