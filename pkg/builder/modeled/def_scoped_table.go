package modeled

import (
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
)

// M builds a [Table] scoped to model type M by scanning a zero value of M.
func M[M Model]() Table[M] {
	return CastT[M](builder.TFrom(new(M)))
}

// CastT wraps a [builder.Table] as a model-scoped [Table].
func CastT[M Model](t builder.Table) Table[M] {
	return &table[M]{Table: t}
}

// Table is a [builder.Table] bound to model type M.
type Table[M Model] interface {
	builder.Table
	ModelNewer[M]

	// MK picks a model-scoped key by name.
	MK(string) Key[M]

	ColIter[M]
	KeyIter[M]
}

type table[M Model] struct {
	Newer[M]
	builder.Table
}

// MCols iterates columns as model-scoped [Col] values.
func (t *table[M]) MCols() iter.Seq[Col[M]] {
	return func(yield func(Col[M]) bool) {
		for c := range t.Cols() {
			if !yield(CastC[M](c)) {
				return
			}
		}
	}
}

// MKeys iterates keys as model-scoped [Key] values.
func (t *table[M]) MKeys() iter.Seq[Key[M]] {
	return func(yield func(Key[M]) bool) {
		for k := range t.Keys() {
			if !yield(CastK[M](k)) {
				return
			}
		}
	}
}

// Unwrap returns the underlying [builder.Table].
func (t *table[M]) Unwrap() builder.Table {
	return t.Table
}

// MK picks a model-scoped key by name.
func (t *table[M]) MK(name string) Key[M] {
	return CastK[M](t.K(name))
}
