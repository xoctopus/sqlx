package internal

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/docx"
	"github.com/xoctopus/x/enumx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/helper"
	"github.com/xoctopus/sqlx/pkg/migrator/models"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
)

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
				Model:   tab.TableName(),
				Col:     c.Name(),
				ColType: def.Type.String(),
				Field:   c.FieldName(),
				Rel:     strings.Join(def.Relation, " "),
				Comment: def.Comment,
			})
		}
	}
	cols, vals := helper.CVsForInsertion(tcs...)

	return append(fragments, frag.Compose("", builder.Insert().Into(t).Values(cols, vals...), frag.Query(";")))
}

func GenerateTableDocuments(_ context.Context, a adaptor.Adaptor, cat builder.Catalog) []frag.Fragment {
	t := builder.TFrom(&models.Table{})
	d := a.Dialect()

	fragments := make([]frag.Fragment, 0)

	fragments = append(fragments, d.DropTable(t))
	fragments = append(fragments, d.CreateTableIfNotExists(t)...)

	tcs := make([]*models.Table, 0)
	for tab := range cat.Tables() {
		v := &models.Table{
			Model:   tab.TableName(),
			TabType: "",
			Comment: "",
		}
		if x, ok := tab.(builder.Newer); ok {
			u := x.New()
			typ := typx.Deref(typx.NewRType(reflect.TypeOf(u)))
			v.TabType = typ.String()
			if p, ok := u.(docx.Provider); ok {
				doc, _ := p.DocOf()
				v.Comment = strings.Join(doc, " ")
			}
		}
		tcs = append(tcs, v)
	}
	cols, vals := helper.CVsForInsertion(tcs...)

	return append(fragments, frag.Compose("", builder.Insert().Into(t).Values(cols, vals...), frag.Query(";")))
}

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
							Model: tab.TableName(),
							Col:   col.Name(),
							Enum:  def.Type.String(),
							Kind:  reflect.TypeOf(e).Kind().String(),
							Key:   e.(interface{ String() string }).String(),
							Value: fmt.Sprintf("%d", e),
							Text:  e.(interface{ Text() string }).Text(),
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
