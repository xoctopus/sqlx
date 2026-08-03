package modeled

import (
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
)

// CastK wraps a [builder.Key] as a model-scoped [Key].
func CastK[M Model](k builder.Key) Key[M] {
	return &key[M]{Key: k}
}

// Key is a [builder.Key] bound to model type M.
type Key[M Model] interface {
	ModelNewer[M]
	builder.Key
	ColIter[M]
}

type key[M Model] struct {
	Newer[M]
	builder.Key
}

// MCols iterates the key's columns as model-scoped [Col] values.
func (k *key[M]) MCols() iter.Seq[Col[M]] {
	return func(yield func(Col[M]) bool) {
		for c := range k.Cols() {
			if !yield(CastC[M](c)) {
				return
			}
		}
	}
}

// KeyIter iterates model-scoped keys.
type KeyIter[M Model] interface {
	MKeys() iter.Seq[Key[M]]
}
