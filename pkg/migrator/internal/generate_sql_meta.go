package internal

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/docx"
	"github.com/xoctopus/x/enumx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/helper"
	"github.com/xoctopus/sqlx/pkg/migrator/models"
)

func cut(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// GenerateTableColumnDocuments rebuilds sql_meta_table_column and inserts column docs from cat.
func GenerateTableColumnDocuments(_ context.Context, a adaptor.Adaptor, cat builder.Catalog) []frag.Fragment {
	t := builder.TFrom(&models.TableColumn{})
	d := a.Dialect()

	fragments := make([]frag.Fragment, 0)

	fragments = append(fragments, d.DropTable(t))
	fragments = append(fragments, d.CreateTableIfNotExists(t)...)

	tcs := make([]*models.TableColumn, 0)
	for tab := range cat.Tables() {
		for c := range tab.Cols() {
			def := builder.GetColDef(c)
			d.ColDefine(&def)
			tcs = append(tcs, &models.TableColumn{
				Model:   cut(tab.TableName(), 64),
				Col:     cut(c.Name(), 64),
				ColType: cut(def.Type.String(), 1024),
				Field:   cut(c.FieldName(), 64),
				Rel:     cut(strings.Join(def.Relation, " "), 128),
				Comment: cut(def.Comment, 1024),
			})
		}
	}
	cols, vals := helper.CVsForInsertion(tcs...)

	return append(fragments, frag.Compose("", builder.Insert().Into(t).Values(cols, vals...), frag.Query(";")))
}

// GenerateTableDocuments rebuilds sql_meta_table and inserts table docs from cat.
func GenerateTableDocuments(_ context.Context, a adaptor.Adaptor, cat builder.Catalog) []frag.Fragment {
	t := builder.TFrom(&models.Table{})
	d := a.Dialect()

	fragments := make([]frag.Fragment, 0)

	fragments = append(fragments, d.DropTable(t))
	fragments = append(fragments, d.CreateTableIfNotExists(t)...)

	tcs := make([]*models.Table, 0)
	for tab := range cat.Tables() {
		v := &models.Table{
			Model:   cut(tab.TableName(), 64),
			TabType: "",
			Comment: "",
		}
		if x, ok := tab.(builder.Newer); ok {
			u := x.New()
			typ := typx.Deref(typx.NewRType(reflect.TypeOf(u)))
			v.TabType = cut(typ.String(), 1024)
			if p, ok := u.(docx.Provider); ok {
				doc, _ := p.DocOf()
				v.Comment = cut(strings.Join(doc, " "), 1024)
			}
		}
		tcs = append(tcs, v)
	}
	cols, vals := helper.CVsForInsertion(tcs...)

	return append(fragments, frag.Compose("", builder.Insert().Into(t).Values(cols, vals...), frag.Query(";")))
}

// GenerateTableEnumerationDocument rebuilds sql_meta_enumeration and inserts enum docs from cat.
func GenerateTableEnumerationDocument(_ context.Context, a adaptor.Adaptor, cat builder.Catalog) []frag.Fragment {
	t := builder.TFrom(&models.Enumeration{})
	d := a.Dialect()

	fragments := make([]frag.Fragment, 0)

	fragments = append(fragments, d.DropTable(t))
	fragments = append(fragments, d.CreateTableIfNotExists(t)...)

	tcs := make([]*models.Enumeration, 0)
	for tab := range cat.Tables() {
		for col := range tab.Cols() {
			def := builder.GetColDef(col)
			if rt, ok := def.Type.Unwrap().(reflect.Type); ok {
				rv := reflect.New(rt)
				if x, ok := rv.Interface().(enumx.CanBeEnum); ok {
					for _, e := range x.EnumValues() {
						v := &models.Enumeration{
							Model:    cut(tab.TableName(), 64),
							Col:      cut(col.Name(), 64),
							EnumType: cut(def.Type.String(), 1024),
							Kind:     cut(reflect.TypeOf(e).Kind().String(), 16),
							Key:      cut(e.(interface{ String() string }).String(), 64),
							Value:    cut(fmt.Sprintf("%d", e), 64),
							Text:     cut(e.(interface{ Text() string }).Text(), 128),
						}
						tcs = append(tcs, v)
					}
				}
			}
		}
	}
	cols, vals := helper.CVsForInsertion(tcs...)

	return append(fragments, frag.Compose("", builder.Insert().Into(t).Values(cols, vals...), frag.Query(";")))
}
