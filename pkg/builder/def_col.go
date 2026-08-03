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
	// Col a database column.
	Col interface {
		frag.Fragment

		// Name returns the database column name.
		Name() string
		// FieldName returns the struct field name of the column definition.
		FieldName() string
		// String returns the qualified name like `table.column`.
		String() string
		// Of rebinds the column to the given table context.
		Of(Table) Col
		// Fragment builds a SQL fragment with `#` as a placeholder for this column.
		Fragment(q string, args ...any) frag.Fragment
	}

	// ColPick picks columns by name.
	ColPick interface {
		// C picks a column by column name or field name.
		C(string) Col
		// Pick returns a column subset by names.
		Pick(...string) Cols
	}

	// ColIter iterates columns.
	ColIter interface {
		// Cols returns an iterator over columns.
		Cols() iter.Seq[Col]
	}

	// ColDef exposes column definition metadata.
	ColDef interface {
		// Def returns the column definition.
		Def() ColumnDef
	}

	// ColWrapper unwraps a wrapped column to its underlying Col.
	ColWrapper interface {
		Unwrap() Col
	}

	// ColModifier mutates column construction options.
	ColModifier interface {
		SetFieldName(string)
		SetComputed(frag.Fragment)
		SetDef(ColumnDef)
	}

	// ColOption configures a column at construction time.
	ColOption func(ColModifier)

	// ColComputed exposes a computed expression for a column.
	ColComputed interface {
		Computed() frag.Fragment
	}

	// ColValuer builds a typed expression against a column.
	ColValuer[T any] func(v Col) frag.Fragment

	// TCol is a typed column that can form conditions and assignments.
	TCol[T any] interface {
		Col

		// AsCond builds a condition with the given valuator.
		AsCond(ColValuer[T]) frag.Fragment
		// AssignBy builds an assignment with the given valuators.
		AssignBy(...ColValuer[T]) Assignment
	}

	// ColsManager adds columns to a collection.
	ColsManager interface {
		AddCol(...Col)
	}

	// Cols is a set of columns.
	Cols interface {
		ColPick
		ColIter

		// Of rebinds all columns to the given table context.
		Of(Table) Cols
		// Len returns the number of columns.
		Len() int

		frag.Fragment
	}

	// ColumnDef is the column definition metadata.
	ColumnDef = def.ColumnDef
)

// C creates a Col by name with optional configuration.
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

// CT creates a typed TCol by name with optional configuration.
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

// CC casts c to a typed TCol[T], optionally applying more options.
func CC[T any](c Col, options ...ColOption) TCol[T] {
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

// GetColTable returns the table bound to c, unwrapping ColWrapper if needed.
func GetColTable(c Col) Table {
	if x, ok := c.(ColWrapper); ok {
		c = x.Unwrap()
	}
	t := Table(nil)
	if x, ok := c.(WithTable); ok {
		t = x.T()
	}
	return t
}

// GetColDef returns the column definition of c, unwrapping ColWrapper if needed.
func GetColDef(c Col) ColumnDef {
	if x, ok := c.(ColWrapper); ok {
		c = x.Unwrap()
	}
	d := ColumnDef{}
	if x, ok := c.(ColDef); ok {
		d = x.Def()
	}
	return d
}

// GetColComputed returns the computed expression of c, unwrapping ColWrapper if needed.
func GetColComputed(c Col) frag.Fragment {
	if x, ok := c.(ColWrapper); ok {
		c = x.Unwrap()
	}
	f := frag.Fragment(nil)
	if x, ok := c.(ColComputed); ok {
		f = x.Computed()
	}
	return f
}

// WithColFieldName sets the struct field name of a column.
func WithColFieldName(name string) ColOption {
	return func(c ColModifier) { c.SetFieldName(name) }
}

// WithColComputed sets a computed expression for a column.
func WithColComputed(f frag.Fragment) ColOption {
	return func(c ColModifier) { c.SetComputed(f) }
}

// WithColDef sets the column definition.
func WithColDef(def *ColumnDef) ColOption {
	return func(c ColModifier) { c.SetDef(*def) }
}

// WithColDefOf parses a ColumnDef from a sample value and struct tag.
func WithColDefOf(v any, tag reflect.StructTag) ColOption {
	return WithColDef(def.ParseColDef(typx.NewRType(reflect.TypeOf(v)), tag))
}

// AsValue assigns/compares using another column as the value.
func AsValue[T any](v TCol[T]) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("?", v)
	}
}

// Value assigns/compares using a literal value.
func Value[T any](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("?", v)
	}
}

// Inc builds `column + v`.
func Inc[T any](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? + ?", c, v)
	}
}

// Dec builds `column - v`.
func Dec[T any](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? - ?", c, v)
	}
}

// Eq builds `column = v`.
func Eq[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? = ?", c, v)
	}
}

// EqCol builds `column = other`.
func EqCol[T comparable](v TCol[T]) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? = ?", c, v)
	}
}

// Neq builds `column <> v`.
func Neq[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? <> ?", c, v)
	}
}

// NeqCol builds `column <> other`.
func NeqCol[T comparable](v TCol[T]) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? <> ?", c, v)
	}
}

// In builds `column IN (...)`. Empty values yield nil.
func In[T any](vs ...T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		if len(vs) == 0 {
			return nil
		}
		return frag.Query("? IN (?)", c, vs)
	}
}

// NotIn builds `column NOT IN (...)`. Empty values yield nil.
func NotIn[T any](vs ...T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		if len(vs) == 0 {
			return nil
		}
		return frag.Query("? NOT IN (?)", c, vs)
	}
}

// IsNull builds `column IS NULL`.
func IsNull[T any]() ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? IS NULL", c)
	}
}

// IsNotNull builds `column IS NOT NULL`.
func IsNotNull[T any]() ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? IS NOT NULL", c)
	}
}

// Like builds `column LIKE %s%`.
func Like[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, "%"+s+"%")
	}
}

// NotLike builds `column NOT LIKE %s%`.
func NotLike[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? NOT LIKE ?", c, "%"+s+"%")
	}
}

// Deprecated: use MatchPrefix or MatchSuffix
func LLike[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, "%"+s)
	}
}

// Deprecated: use MatchPrefix or MatchSuffix
func RLike[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, s+"%")
	}
}

// MatchPrefix builds `column LIKE s%`.
func MatchPrefix[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, s+"%")
	}
}

// MatchSuffix builds `column LIKE %s`.
func MatchSuffix[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, "%"+s)
	}
}

// Match builds `column LIKE s` without rewriting the pattern.
func Match[T ~string](s T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? LIKE ?", c, s)
	}
}

// Between builds `column BETWEEN min AND max`.
func Between[T comparable](min, max T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? BETWEEN ? AND ?", c, min, max)
	}
}

// NotBetween builds `column NOT BETWEEN min AND max`.
func NotBetween[T comparable](min, max T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? NOT BETWEEN ? AND ?", c, min, max)
	}
}

// Gt builds `column > v`.
func Gt[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? > ?", c, v)
	}
}

// Gte builds `column >= v`.
func Gte[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? >= ?", c, v)
	}
}

// Lt builds `column < v`.
func Lt[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? < ?", c, v)
	}
}

// Lte builds `column <= v`.
func Lte[T comparable](v T) ColValuer[T] {
	return func(c Col) frag.Fragment {
		return frag.Query("? <= ?", c, v)
	}
}

// Columns creates a Cols from column names.
func Columns(names ...string) Cols {
	cs := &columns{}
	for _, name := range names {
		cs.AddCol(C(name))
	}
	return cs
}

// ColsOf creates a Cols from the given columns.
func ColsOf(cs ...Col) Cols {
	cs_ := &columns{}
	for _, c := range cs {
		cs_.AddCol(c)
	}
	return cs_
}

// ColsIterOf collects columns from one or more iterators into a Cols.
func ColsIterOf(cs ...iter.Seq[Col]) Cols {
	cs_ := &columns{}
	for _, seq := range cs {
		for c := range seq {
			cs_.AddCol(c)
		}
	}
	return cs_
}

type columns struct {
	l []Col
}

// C picks a Col by field name (exported) or column name.
func (cs *columns) C(name string) Col {
	if name == "" {
		return nil
	}

	for i := range cs.l {
		c := cs.l[i]
		if x := name[0]; x >= 'A' && x <= 'Z' && c.FieldName() == name {
			return c
		}

		if c.Name() == name {
			return c
		}
	}
	return nil
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
			if !yield(c) {
				return
			}
		}
	}
}

func (cs *columns) Pick(names ...string) Cols {
	sub := &columns{}
	if len(names) == 0 {
		return cs
	}
	for _, name := range names {
		c := cs.C(name)
		must.NotNilF(c, "unknown column %s from %v", name, names)
		sub.AddCol(c)
	}
	return sub
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
				if !yield(c) {
					return
				}
			}
		}

		for q, args := range frag.ComposeSeq(",", fragments).Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}
	}
}
