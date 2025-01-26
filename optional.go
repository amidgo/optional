package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var ErrOptionalValueIsEmpty = errors.New("optional value is empty")

type (
	Byte = Optional[byte]

	Uint   = Optional[uint]
	Uint8  = Optional[uint8]
	Uint16 = Optional[uint16]
	Uint32 = Optional[uint32]
	Uint64 = Optional[uint64]

	Int   = Optional[int]
	Int8  = Optional[int8]
	Int16 = Optional[int16]
	Int32 = Optional[int32]
	Int64 = Optional[int64]

	Bool = Optional[bool]

	String = Optional[string]

	Float32 = Optional[float32]
	Float64 = Optional[float64]
)

type Optional[T any] struct {
	ok    bool
	value T
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
