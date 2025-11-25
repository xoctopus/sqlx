package modeled

import (
	"context"
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
)

func M[M Model](ctx context.Context) Table[M] {
	return CastT[M](builder.TFrom(ctx, new(M)))
}

func CastT[M Model](t builder.Table) Table[M] {
	return &table[M]{Table: t}
}

type Table[M Model] interface {
	builder.Table
	ModelNewer[M]
	ColIter[M]
	KeyIter[M]
}

type table[M Model] struct {
	Newer[M]
	builder.Table
}

func (t *table[M]) MCols() iter.Seq[Col[M]] {
	return func(yield func(Col[M]) bool) {
		for c := range t.Cols() {
			yield(CastC[M](c))
		}
	}
}

func (t *table[M]) MKeys() iter.Seq[Key[M]] {
	return func(yield func(Key[M]) bool) {
		for k := range t.Keys() {
			yield(CastK[M](k))
		}
	}
}

func (t *table[m]) Unwrap() builder.Table {
	return t.Table
}
