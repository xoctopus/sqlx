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
	Key interface {
		frag.Fragment
		ColIter

		Name() string
		Of(Table) Key
		IsPrimary() bool
		IsUnique() bool
		// String return [$table_name.]$index_name
		String() string
	}

	KeyDefine = def.KeyDefine

	KeyDef interface {
		Method() string
		ColumnOptions() []def.KeyColumnOption
	}

	KeyPick interface {
		K(string) Key
	}

	KeyIter interface {
		Keys() iter.Seq[Key]
	}

	KeysManager interface {
		AddKey(...Key)
	}

	Keys interface {
		KeyIter
		KeyPick

		Of(Table) Keys
		// Len() int
	}

	KeyKind         = def.KeyKind
	KeyColumnOption = def.KeyColumnOption
)

func PK(cols Cols, opts ...KeyOption) Key {
	must.BeTrueF(cols != nil && cols.Len() > 0, "missing columns to create primary key")
	return UK("PRIMARY", cols, opts...)
}

func UK(name string, cols Cols, opts ...KeyOption) Key {
	must.BeTrueF(cols != nil && cols.Len() > 0, "missing columns to create unique index")
	return K(name, cols, append(opts, WithKeyUniqueness(true))...)
}

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

type KeyOption func(*key)

func WithKeyUniqueness(unique bool) KeyOption {
	return func(k *key) {
		k.unique = unique
	}
}

func WithKeyMethod(method string) KeyOption {
	return func(k *key) {
		k.method = method
	}
}

func WithKeyColumnOptions(opts ...KeyColumnOption) KeyOption {
	return func(k *key) {
		k.options = opts
	}
}

func GetKeyTable(k Key) Table {
	if d, ok := k.(WithTable); ok {
		return d.T()
	}
	return nil
}

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
					yield(",", nil)
				}
				c := cols.C(o.Name)
				must.BeTrueF(c != nil, "missing column: %s", o.Name)
				yield(c.Name(), nil)
				if len(o.Options) > 0 {
					yield(" "+strings.Join(o.Options, " "), nil)
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
		names := map[string]bool{}
		for _, opt := range k.options {
			names[opt.Name] = true
		}
		for c := range k.table.Cols() {
			if names[c.FieldName()] || names[c.Name()] {
				yield(c)
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
			yield(k)
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
