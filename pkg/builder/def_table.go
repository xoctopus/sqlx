package builder

import (
	"context"
	"iter"
	"strings"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type (
	Table interface {
		TableName() string
		String() string

		KeyPick
		KeyIter

		ColPick
		ColIter

		Fragment(string, ...any) frag.Fragment
		frag.Fragment
	}

	Tables interface {
		Table(string) Table
		Tables() iter.Seq[Table]

		Add(...Table)
		Remove(...Table)

		Require(...Table)
	}

	WithTable interface {
		T() Table
	}

	WithTableName interface {
		WithTableName(name string) Table
	}
)

func T(name string, defs ...frag.Fragment) Table {
	t := &table{
		name: name,
		cs:   &columns{},
		ks:   &keys{},
	}

	// col added first
	for _, def := range defs {
		if c, ok := def.(Col); ok {
			t.AddCol(c)
		}
	}

	for _, def := range defs {
		if k, ok := def.(Key); ok {
			t.AddKey(k)
		}
	}

	return t
}

type table struct {
	database string
	name     string
	desc     []string

	ks Keys
	cs Cols
}

func (t *table) AddCol(cols ...Col) {
	for _, c := range cols {
		t.cs.(ColsManager).AddCol(c.Of(t))
	}
}

func (t *table) AddKey(keys ...Key) {
	for _, k := range keys {
		t.ks.(KeysManager).AddKey(k.Of(t))
	}
}

func (t *table) TableName() string {
	return t.name
}

func (t *table) WithTableName(name string) Table {
	tt := &table{
		name:     name,
		database: t.database,
		desc:     t.desc,
	}

	tt.cs = t.cs.Of(tt)
	tt.ks = t.ks.Of(tt)

	return tt
}

func (t *table) String() string {
	return t.name
}

func (t *table) IsNil() bool {
	return t == nil || t.name == ""
}

func (t *table) Frag(ctx context.Context) frag.Iter {
	return frag.Lit(t.name).Frag(ctx)
}

func (t *table) Fragment(query string, args ...any) frag.Fragment {
	if query == "" {
		return nil
	}

	set := frag.NamedArgs{
		"_t_": t,
	}

	for col := range t.cs.Cols() {
		set["_t_"+col.FieldName()] = col
	}

	q := strings.ReplaceAll(query, "#", "@_t_")

	return frag.Query(q, append([]any{set}, args...)...)
}

func (t *table) C(name string) Col {
	return t.cs.C(name)
}

func (t *table) Cols() iter.Seq[Col] {
	return t.cs.Cols()
}

func (t *table) K(name string) Key {
	return t.ks.K(name)
}

func (t *table) Keys() iter.Seq[Key] {
	return t.ks.Keys()
}
