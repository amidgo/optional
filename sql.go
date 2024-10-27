package optional

import "database/sql"

func Null[T any](op Optional[T]) sql.Null[T] {
	return sql.Null[T]{
		V:     op.value,
		Valid: op.ok,
	}
}
