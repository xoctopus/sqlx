package modeled

import (
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
)

func CastK[M Model](k builder.Key) Key[M] {
	return &key[M]{Key: k}
}

type Key[M Model] interface {
	ModelNewer[M]
	builder.Key
	ColIter[M]
}

type key[M Model] struct {
	Newer[M]
	builder.Key
}

func (k *key[M]) MCols() iter.Seq[Col[M]] {
	return func(yield func(Col[M]) bool) {
		for c := range k.Cols() {
			if !yield(CastC[M](c)) {
				return
			}
		}
	}
}

type KeyIter[M Model] interface {
	MKeys() iter.Seq[Key[M]]
}
