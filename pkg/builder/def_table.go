package builder

import (
	"container/list"
	"context"
	"iter"
	"reflect"
	"strings"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/internal"
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

	Catalog interface {
		T(string) Table
		Tables() iter.Seq[Table]

		Add(...Table)
		Remove(string)
		Len() int

		// Require(...Catalog)
	}

	WithTable interface {
		T() Table
	}
	WithTableName interface {
		WithTableName(name string) Table
	}
	WithSchema interface {
		WithSchema(schema string) Table
	}
	HasSchema interface {
		Schema() string
	}
	WithDatabase interface {
		WithDatabase(database string) Table
	}
	HasDatabase interface {
		Database() string
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

var schemas = syncx.NewXmap[reflect.Type, Table]()

func TFrom(m any) Table {
	t := reflect.TypeOf(m)
	must.BeTrueF(t.Kind() == reflect.Pointer, "model %s must be a pointer", t.Name())
	t = t.Elem()
	must.BeTrueF(t.Kind() == reflect.Struct, "model %s must be a struct", t.Name())

	if tab, ok := schemas.Load(t); ok {
		if x, ok := m.(internal.Model); ok {
			return tab.(WithTableName).WithTableName(x.TableName())
		}
		return tab
	}

	name := t.Name()
	if x, ok := m.(internal.Model); ok {
		name = x.TableName()
	}

	tab := scan(m)
	tab = tab.(WithTableName).WithTableName(name)
	schemas.Store(t, tab)

	return tab
}

type table struct {
	// database identifies a unique endpoint
	database string
	// schema identifies a unique schema under endpoint
	schema string

	name string
	desc []string

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

func (t *table) WithSchema(schema string) Table {
	t_ := *t
	t_.schema = schema
	return &t_
}

func (t *table) Schema() string {
	return t.schema
}

func (t *table) WithDatabase(database string) Table {
	t_ := *t
	t_.database = database
	return &t_
}

func (t *table) Database() string {
	return t.database
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

func (t *table) Pick(names ...string) Cols {
	return t.cs.Pick(names...)
}

func (t *table) K(name string) Key {
	return t.ks.K(name)
}

func (t *table) Keys() iter.Seq[Key] {
	return t.ks.Keys()
}

func TableNames(c Catalog) iter.Seq[string] {
	return func(yield func(string) bool) {
		for t := range c.Tables() {
			yield(t.TableName())
		}
	}
}

func CatalogFrom(models ...internal.Model) Catalog {
	tabs := &tables{}
	for i := range models {
		tabs.Add(TFrom(models[i]))
	}
	return tabs
}

func NewCatalog() Catalog {
	return &tables{}
}

type tables struct {
	l *list.List
	m map[string]*list.Element
	// requires []Catalog
}

func (t *tables) T(name string) Table {
	if t.m != nil {
		if x, ok := t.m[name]; ok {
			return x.Value.(Table)
		}
	}
	return nil
}

func (t *tables) Tables() iter.Seq[Table] {
	return func(yield func(Table) bool) {
		// emitted := make(map[string]bool)

		// emit := func(t Table) bool {
		// 	name := t.TableName()
		// 	if _, ok := emitted[name]; ok {
		// 		return true
		// 	}
		// 	emitted[name] = true
		// 	return yield(t)
		// }

		if t.l != nil {
			for e := t.l.Front(); e != nil; e = e.Next() {
				// x := e.Value.(Table)
				// emit(x)
				yield(e.Value.(Table))
			}
		}

		// for _, c := range t.requires {
		// 	for x := range c.Tables() {
		// 		emit(x)
		// 	}
		// }
	}
}

func (t *tables) Add(tables ...Table) {
	if t.m == nil {
		t.m = make(map[string]*list.Element)
		t.l = list.New()
	}

	for _, x := range tables {
		if x != nil {
			name := x.TableName()
			if _, ok := t.m[name]; ok {
				t.Remove(name)
			}
			t.m[name] = t.l.PushBack(x)
		}
	}
}

func (t *tables) Remove(name string) {
	if t.m != nil {
		if e, ok := t.m[name]; ok {
			t.l.Remove(e)
			delete(t.m, name)
		}
	}
}

func (t *tables) Len() int {
	if t.m != nil {
		return len(t.m)
	}
	return 0
}

// func (t *tables) Require(requires ...Catalog) {
// 	t.requires = append(t.requires, requires...)
// }
