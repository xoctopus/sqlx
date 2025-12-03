package builder

import (
	"context"
	"database/sql"
	"iter"
	"reflect"
	"strings"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type (
	Col interface {
		frag.Fragment

		// Name presents database column name
		Name() string
		// FieldName presents name of definition of struct field
		FieldName() string
		// String returns full name like `table.column`
		String() string
		// Of change column table context to given argument
		Of(Table) Col
		// Fragment return a sql frag based current column
		Fragment(q string, args ...any) frag.Fragment
	}

	ColPick interface {
		// C picks a column by name
		C(string) Col
	}

	ColIter interface {
		// Cols iteration of columns
		Cols() iter.Seq[Col]
	}

	ColDef interface {
		// Def returns column define
		Def() ColumnDef
	}

	ColWrapper interface {
		Unwrap() Col
	}

	ColModifier interface {
		SetFieldName(string)
		SetComputed(frag.Fragment)
		SetDef(ColumnDef)
	}

	ColOption func(ColModifier)

	ColComputed interface {
		Computed() frag.Fragment
	}

	ColValuer[T any] func(v Col) frag.Fragment

	// TCol make column can computed
	TCol[T any] interface {
		Col

		AsCond(ColValuer[T]) frag.Fragment
		AssignBy(...ColValuer[T]) Assignment
	}

	ColsManager interface {
		AddCol(...Col)
	}

	Cols interface {
		ColPick
		ColIter

		// Of changes columns table context
		Of(Table) Cols
		// Len returns amount of column set
		Len() int

		frag.Fragment
	}

	ColumnDef = def.ColumnDef
)

// C creates a Col by name and ColOption
func C(name string, options ...ColOption) Col {
	c := &column[any]{
		name: strings.ToLower(name),
		def:  ColumnDef{},
	}

	for _, o := range options {
		o(c)
	}
	return c
}

// CT creates a Col by name and ColOption with computing T
func CT[T any](name string, options ...ColOption) TCol[T] {
	c := &column[T]{
		name: strings.ToLower(name),
		def:  ColumnDef{},
	}
	for _, apply := range options {
		apply(c)
	}
	return c
}

func CastC[T any](c Col, options ...ColOption) TCol[T] {
	col := &column[T]{
		name:     c.Name(),
		fname:    c.FieldName(),
		table:    GetColTable(c),
		def:      GetColDef(c),
		computed: GetColComputed(c),
	}

	for _, o := range options {
		o(col)
	}
	return col
}

func PickCols(p ColPick, names ...string) Cols {
	cs := &columns{}
	for _, name := range names {
		c := p.C(name)
		must.NotNilF(c, "unknown column %s from %v", name, names)
		cs.AddCol(c)
	}
	return cs
}

type column[T any] struct {
	name     string
	fname    string
	def      ColumnDef
	table    Table
	computed frag.Fragment
}

func (c *column[T]) FieldName() string {
	return c.fname
}

func (c *column[T]) SetFieldName(name string) {
	c.fname = name
}

func (c *column[T]) Computed() frag.Fragment {
	return c.computed
}

func (c *column[T]) SetComputed(f frag.Fragment) {
	c.computed = f
}

func (c *column[T]) Def() ColumnDef {
	return c.def
}

func (c *column[T]) SetDef(def ColumnDef) {
	c.def = def
}

func (c *column[T]) Name() string {
	return c.name
}

func (c *column[T]) T() Table {
	return c.table
}

func (c *column[T]) Of(t Table) Col {
	return &column[T]{
		name:     c.name,
		fname:    c.fname,
		def:      c.def,
		computed: c.computed,
		table:    t,
	}
}

func (c *column[T]) String() string {
	if c.table != nil {
		return c.table.TableName() + "." + c.name
	}
	return c.name
}

func (c *column[T]) Fragment(q string, args ...any) frag.Fragment {
	q = strings.ReplaceAll(q, "#", "@_column")
	return frag.Query(q, append([]any{sql.Named("_column", c)}, args)...)
}

func (c *column[T]) AsCond(op ColValuer[T]) frag.Fragment {
	if op != nil {
		return op(c)
	}
	return nil
}

func (c *column[T]) AssignBy(ops ...ColValuer[T]) Assignment {
	if len(ops) == 0 {
		return nil
	}
	vs := make([]any, 0, len(ops))
	for _, op := range ops {
		if op != nil {
			vs = append(vs, op(c))
		}
	}
	return ColumnsAndValues(c, vs...)
}

func (c *column[T]) IsNil() bool {
	return c == nil
}

func (c *column[T]) Frag(ctx context.Context) frag.Iter {
	toggles := TogglesFromContext(ctx)

	if c.computed != nil && toggles.Is(TOGGLE__IN_PROJECT) {
		return frag.Query("? AS ?", c.computed, frag.Lit(c.name)).Frag(ctx)
	}

	if toggles.Is(TOGGLE__MULTI_TABLE) {
		must.BeTrueF(c.table != nil, "table is not define on column: %s", c.name)
		if toggles.Is(TOGGLE__AUTO_ALIAS) {
			return frag.Query(
				"?.? AS ?",
				c.table,
				frag.Lit(c.name),
				frag.Query(frag.Alias(c.table.TableName(), c.name)),
			).Frag(ctx)
		}
		return frag.Query("?.?", c.table, frag.Lit(c.name)).Frag(ctx)
	}
	return frag.Lit(c.name).Frag(ctx)
}

func GetColTable(c Col) Table {
	if x, ok := c.(ColWrapper); ok {
		c = x.Unwrap()
	}
	if x, ok := c.(WithTable); ok {
		return x.T()
	}
	return nil
}

func GetColDef(c Col) ColumnDef {
	if x, ok := c.(ColWrapper); ok {
		c = x.Unwrap()
	}
	if x, ok := c.(ColDef); ok {
		return x.Def()
	}
	return ColumnDef{}
}

func GetColComputed(c Col) frag.Fragment {
	if x, ok := c.(ColWrapper); ok {
		c = x.Unwrap()
	}
	if x, ok := c.(ColComputed); ok {
		return x.Computed()
	}
	return nil
}

func WithColFieldName(name string) ColOption {
	return func(c ColModifier) { c.SetFieldName(name) }
}

func WithColComputed(f frag.Fragment) ColOption {
	return func(c ColModifier) { c.SetComputed(f) }
}

func WithColDef(def *ColumnDef) ColOption {
	return func(c ColModifier) { c.SetDef(*def) }
}

func WithColDefOf(ctx context.Context, v any, tag reflect.StructTag) ColOption {
	return WithColDef(def.ParseColDef(ctx, typx.NewRType(reflect.TypeOf(v)), tag))
}

func AsValue[T any](v TCol[T]) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("?", v)
	}
}

func Value[T any](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("?", v)
	}
}

func Inc[T any](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? + ?", c, v)
	}
}

func Dec[T any](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? - ?", c, v)
	}
}

func Eq[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? = ?", c, v)
	}
}

func EqCol[T comparable](v TCol[T]) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? = ?", c, v)
	}
}

func Neq[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? <> ?", c, v)
	}
}

func NeqCol[T comparable](v TCol[T]) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? <> ?", c, v)
	}
}

func In[T any](vs ...T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		if len(vs) == 0 {
			return nil
		}
		return frag.Query("? IN (?)", c, vs)
	}
}

func NotIn[T any](vs ...T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		if len(vs) == 0 {
			return nil
		}
		return frag.Query("? NOT IN (?)", c, vs)
	}
}

func IsNull[T any]() ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? IS NULL", c)
	}
}

func IsNotNull[T any]() ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? IS NOT NULL", c)
	}
}

func Like[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, "%"+s+"%")
	}
}

func NotLike[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? NOT LIKE ?", c, "%"+s+"%")
	}
}

func LLike[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, "%"+s)
	}
}

func RLike[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, s+"%")
	}
}

func Between[T comparable](min, max T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? BETWEEN ? AND ?", c, min, max)
	}
}

func NotBetween[T comparable](min, max T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? NOT BETWEEN ? AND ?", c, min, max)
	}
}

func Gt[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? > ?", c, v)
	}
}

func Gte[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? >= ?", c, v)
	}
}

func Lt[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? < ?", c, v)
	}
}

func Lte[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? <= ?", c, v)
	}
}

func Columns(names ...string) Cols {
	cs := &columns{}
	for _, name := range names {
		cs.AddCol(C(name))
	}
	return cs
}

func ColsOf(cs ...Col) Cols {
	cs_ := &columns{}
	for _, c := range cs {
		cs_.AddCol(c)
	}
	return cs_
}

type columns struct {
	l []Col
}

func (cs *columns) F(name string) Col {
	for i := range cs.l {
		c := cs.l[i]
		if MatchColumn(c, name) {
			return c
		}
	}
	return nil
}

// C collects Col from columns by column name or field name
func (cs *columns) C(name string) Col {
	return cs.F(name)
}

func (cs *columns) Len() int {
	if cs == nil || cs.l == nil {
		return 0
	}
	return len(cs.l)
}

func (cs *columns) Cols() iter.Seq[Col] {
	return func(yield func(Col) bool) {
		for _, c := range cs.l {
			yield(c)
		}
	}
}

func (cs *columns) AddCol(cols ...Col) {
	for i := range cols {
		if c := cols[i]; c != nil {
			cs.l = append(cs.l, c)
		}
	}
}

func (cs *columns) Of(t Table) Cols {
	cs2 := &columns{}

	for i := range cs.l {
		cs2.AddCol(cs.l[i].Of(t))
	}

	return cs2
}

func (cs *columns) IsNil() bool {
	return cs == nil || cs.Len() == 0
}

func (cs *columns) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		fragments := func(yield func(frag.Fragment) bool) {
			for _, c := range cs.l {
				yield(c)
			}
		}

		for q, args := range frag.ComposeSeq(",", fragments).Frag(ctx) {
			yield(q, args)
		}
	}
}
