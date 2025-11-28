package structs

import (
	"context"
	"database/sql/driver"
	"go/ast"
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/xoctopus/typex"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/internal"
	"github.com/xoctopus/sqlx/internal/def"
)

var (
	cache         = syncx.NewXmap[any, []*Field]()
	tSqlModel     = reflect.TypeFor[internal.Model]()
	tDriverValuer = reflect.TypeFor[driver.Valuer]()
)

func FieldsFor(ctx context.Context, t typex.Type) []*Field {
	t = typex.Deref(t)
	must.BeTrueF(
		t.Kind() == reflect.Struct,
		"model %s must be a struct, but got %s",
		t.Name(), t.Kind(),
	)

	fields, ok := cache.Load(t.Unwrap())
	if ok {
		return fields
	}

	defer func() {
		cache.Store(t.Unwrap(), fields)
	}()

	for f := range (&walker{}).Walk(ctx, t) {
		fields = append(fields, f)
	}
	return fields
}

func FieldsSeqFor(ctx context.Context, t typex.Type) iter.Seq[*Field] {
	return slices.Values(FieldsFor(ctx, t))
}

type walker struct {
	flocs []int
	mlocs []int
	t     typex.Type
}

func (w *walker) Walk(ctx context.Context, t typex.Type) iter.Seq[*Field] {
	mlocs := w.mlocs[:]
	mtype := w.t

	if ok := t.Implements(typex.NewRType(ctx, tSqlModel)); ok {
		if mtype != nil && mtype.NumField() == 1 && mtype.Field(0).Anonymous() {
			// extendable
		} else {
			mtype = t
			mlocs = w.flocs[:]
		}
	}

	return func(yield func(*Field) bool) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !ast.IsExported(f.Name()) {
				continue
			}

			tag := def.ModelTagKeyFrom(ctx)
			loc := append(w.flocs, i)
			flag := reflectx.ParseTag(f.Tag()).Get(tag)
			name := f.Name()
			if flag != nil {
				if flag.Name() == "-" {
					continue
				}
				if flag.Name() != "" {
					name = flag.Name()
				}
			}

			if (f.Anonymous() || f.Type().Name() == f.Name()) && flag == nil {
				ft := f.Type()

				if !ft.Implements(typex.NewRType(ctx, tDriverValuer)) {
					ft = typex.Deref(ft)
					if ft.Kind() == reflect.Struct {
						embed := &walker{
							flocs: loc,
							mlocs: mlocs,
							t:     mtype,
						}
						for c := range embed.Walk(ctx, ft) {
							if !yield(c) {
								return
							}
						}
						continue
					}
				}
			}

			p := &Field{
				Name:      strings.ToLower(name),
				FieldName: f.Name(),
				Type:      f.Type(),
				Field:     f,
				Flag:      flag,
				ColumnDef: *def.ParseColDef(ctx, f.Type(), f.Tag()),
			}
			p.Loc = make([]int, len(loc))
			copy(p.Loc, loc)
			// p.ModelLoc = make([]int, len(mlocs))
			// copy(p.ModelLoc, mlocs)

			if !yield(p) {
				return
			}
		}
	}
}

type Field struct {
	Name      string
	FieldName string
	Type      typex.Type
	Field     typex.StructField
	Flag      *reflectx.Flag
	Loc       []int
	ColumnDef def.ColumnDef
	// ModelLoc  []int
}

func (f *Field) Value(v reflect.Value) any {
	return f.FieldValue(v).Interface()
}

func (f *Field) FieldValue(v reflect.Value) reflect.Value {
	return value(v, f.Loc)
}

func (f *Field) ModelValue() reflect.Value {
	panic("forbidden")
}

func value(v reflect.Value, indexes []int) reflect.Value {
	must.BeTrueF(v.CanSet(), "struct value must be able to set")
	n := len(indexes)
	if n == 0 {
		return v
	}
	v = reflectx.Indirect(v)
	fv := v

	for i := 0; i < n; i++ {
		idx := indexes[i]
		fv = fv.Field(idx)

		// last loc should keep ptr value
		if i < n-1 {
			for fv.Kind() == reflect.Pointer {
				// notice the ptr struct ensure only for Ptr Anonymous Field
				if fv.IsNil() {
					fv.Set(reflectx.New(fv.Type()))
				}
				fv = fv.Elem()
			}
		}
	}
	return fv
}

func TableFields(ctx context.Context, v any) []*TableField {
	return slices.Collect(TableFieldsSeq(ctx, v))
}

func TableFieldsSeq(ctx context.Context, v any) iter.Seq[*TableField] {
	return func(yield func(*TableField) bool) {
		rv := reflectx.Indirect(reflect.ValueOf(v))
		must.BeTrueF(rv.IsValid(), "struct value must be valid")

		tableName := ""
		if m, ok := v.(internal.Model); ok {
			tableName = m.TableName()
		}

		for f := range FieldsSeqFor(ctx, typex.NewRType(ctx, rv.Type())) {
			if f.Flag != nil && f.Flag.Option("deprecated") != nil {
				continue
			}
			tf := &TableField{
				Field:     *f,
				Value:     f.FieldValue(rv),
				TableName: tableName,
			}
			if !yield(tf) {
				return
			}
		}
	}
}

type TableField struct {
	Field     Field
	Value     reflect.Value
	TableName string
}
