package builder

import (
	"context"
	"iter"
	"strings"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type (
	// Key is a table index or primary key.
	Key interface {
		frag.Fragment
		ColIter

		// Name returns the key/index name.
		Name() string
		// Of rebinds the key to the given table context.
		Of(Table) Key
		// IsPrimary reports whether this is a primary key.
		IsPrimary() bool
		// IsUnique reports whether this is a unique key.
		IsUnique() bool
		// String returns [$table_name.]$index_name.
		String() string
	}

	// KeyDefine is the key definition metadata.
	KeyDefine = def.KeyDefine

	// KeyDef exposes index method and per-column options.
	KeyDef interface {
		Method() string
		ColumnOptions() []def.KeyColumnOption
	}

	// KeyPick picks a key by name.
	KeyPick interface {
		K(string) Key
	}

	// KeyIter iterates keys.
	KeyIter interface {
		Keys() iter.Seq[Key]
	}

	// KeysManager adds keys to a collection.
	KeysManager interface {
		AddKey(...Key)
	}

	// Keys is a set of keys.
	Keys interface {
		KeyIter
		KeyPick

		// Of rebinds all keys to the given table context.
		Of(Table) Keys
		// Len() int
	}

	// KeyKind is the kind of a key.
	KeyKind = def.KeyKind
	// KeyColumnOption is a per-column key option.
	KeyColumnOption = def.KeyColumnOption
)

// PK creates a primary key over cols.
func PK(cols Cols, opts ...KeyOption) Key {
	must.BeTrueF(cols != nil && cols.Len() > 0, "missing columns to create primary key")
	return UK("PRIMARY", cols, opts...)
}

// UK creates a unique index named name over cols.
func UK(name string, cols Cols, opts ...KeyOption) Key {
	must.BeTrueF(cols != nil && cols.Len() > 0, "missing columns to create unique index")
	return K(name, cols, append(opts, WithKeyUniqueness(true))...)
}

// K creates an index named name over cols.
func K(name string, cols Cols, opts ...KeyOption) Key {
	must.BeTrueF(cols != nil && cols.Len() > 0, "missing columns to create index")
	k := &key{name: strings.ToLower(name)}

	for c := range cols.Cols() {
		k.options = append(k.options, def.KeyColumnOption{Name: c.Name()})
	}

	for _, f := range opts {
		f(k)
	}
	return k
}

// KeyOption configures a key at construction time.
type KeyOption func(*key)

// WithKeyUniqueness marks the key as unique or non-unique.
func WithKeyUniqueness(unique bool) KeyOption {
	return func(k *key) {
		k.unique = unique
	}
}

// WithKeyMethod sets the index method (for example BTREE).
func WithKeyMethod(method string) KeyOption {
	return func(k *key) {
		k.method = method
	}
}

// WithKeyColumnOptions sets per-column key options.
func WithKeyColumnOptions(opts ...KeyColumnOption) KeyOption {
	return func(k *key) {
		options := make([]KeyColumnOption, 0)
		for _, o := range opts {
			if len(o.Name) > 0 {
				options = append(options, o)
			}
		}
		k.options = options
	}
}

// GetKeyTable returns the table bound to k, if any.
func GetKeyTable(k Key) Table {
	if d, ok := k.(WithTable); ok {
		return d.T()
	}
	return nil
}

// KeyColumnsDefOf renders key column definitions for DDL.
func KeyColumnsDefOf(k Key) frag.Fragment {
	kd := k.(KeyDef)

	must.BeTrueF(
		len(kd.ColumnOptions()) > 0,
		"missing columns of key define: %s", k,
	)

	cols := ColsIterOf(k.Cols())
	return frag.Func(func(ctx context.Context) frag.Iter {
		return func(yield func(string, []any) bool) {
			for i, o := range kd.ColumnOptions() {
				if i > 0 {
					if !yield(",", nil) {
						return
					}
				}
				c := cols.C(o.Name)
				must.BeTrueF(c != nil, "missing column: %s", o.Name)
				if !yield(c.Name(), nil) {
					return
				}
				if len(o.Options) > 0 {
					if !yield(" "+strings.Join(o.Options, " "), nil) {
						return
					}
				}
			}
		}
	})
}

type key struct {
	table   Table
	kind    KeyKind
	name    string
	unique  bool
	method  string
	options []def.KeyColumnOption
}

func (k *key) IsNil() bool { return k == nil }

func (k *key) Frag(ctx context.Context) frag.Iter {
	return frag.Lit(k.name).Frag(ctx)
}

func (k *key) T() Table {
	return k.table
}

func (k *key) Method() string {
	return k.method
}

func (k *key) ColumnOptions() []def.KeyColumnOption {
	return k.options
}

func (k *key) Name() string {
	return k.name
}

func (k *key) IsUnique() bool {
	return k.unique
}

func (k *key) IsPrimary() bool {
	return k.unique && (k.name == "primary" || strings.HasSuffix(k.name, "pkey"))
}

func (k *key) String() string {
	s := ""
	if k.table != nil {
		s += k.table.String() + "_"
	}
	return s + k.name
}

func (k *key) Cols() iter.Seq[Col] {
	return func(yield func(Col) bool) {
		if k.table == nil {
			return
		}
		names := map[string]bool{}
		for _, opt := range k.options {
			names[opt.Name] = true
		}
		for c := range k.table.Cols() {
			if names[c.FieldName()] || names[c.Name()] {
				if !yield(c) {
					return
				}
			}
		}
	}
}

func (k *key) Of(t Table) Key {
	k_ := *k
	k_.table = t
	return &k_
}

type keys struct {
	ks []Key
}

func (ks *keys) K(name string) Key {
	name = strings.ToLower(name)
	for i := range ks.ks {
		if name == ks.ks[i].Name() {
			return ks.ks[i]
		}
	}
	return nil
}

func (ks *keys) AddKey(followers ...Key) {
	for i := range followers {
		k := followers[i]
		if k != nil {
			ks.ks = append(ks.ks, k)
		}
	}
}

func (ks *keys) Keys() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		for _, k := range ks.ks {
			if !yield(k) {
				return
			}
		}
	}
}

func (ks *keys) Of(t Table) Keys {
	cloned := &keys{}
	for i := range ks.ks {
		cloned.AddKey(ks.ks[i].Of(t))
	}
	return cloned
}
