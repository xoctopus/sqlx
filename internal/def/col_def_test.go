package def_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/xoctopus/typex"
	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
)

func TestParseColDef(t *testing.T) {
	typ := typex.NewRType(context.Background(), reflect.TypeFor[int]())
	var cases = []struct {
		name string
		def  *def.ColumnDef
		tag  reflect.StructTag
	}{
		{
			name: "AutoInc",
			def:  &def.ColumnDef{AutoInc: true},
			tag:  reflect.StructTag(`db:",autoinc"`),
		},
		{
			name: "Null",
			def:  &def.ColumnDef{Null: true},
			tag:  reflect.StructTag(`db:",null"`),
		},
		{
			name: "Width",
			def:  &def.ColumnDef{Width: 10},
			tag:  reflect.StructTag(`db:",width=10"`),
		},
		{
			name: "Precision",
			def:  &def.ColumnDef{Precision: 10},
			tag:  reflect.StructTag(`db:",precision=10"`),
		},
		{
			name: "Default",
			def:  &def.ColumnDef{Default: ptrx.Ptr("abc def")},
			tag:  reflect.StructTag(`db:",default='abc def'"`),
		},
		{
			name: "DefaultNull",
			def:  &def.ColumnDef{Default: ptrx.Ptr("")},
			tag:  reflect.StructTag(`db:",default=''"`),
		},
		{
			name: "OnUpdate",
			def:  &def.ColumnDef{OnUpdate: ptrx.Ptr("CURRENT_TIMESTAMP")},
			tag:  reflect.StructTag(`db:",onupdate='CURRENT_TIMESTAMP'"`),
		},
		{
			name: "Deprecated",
			def:  &def.ColumnDef{Deprecated: &def.DeprecatedActions{RenameTo: "f_new_column_name"}},
			tag:  reflect.StructTag(`db:",deprecated='f_new_column_name'"`),
		},
		{
			name: "NoFlag",
			def:  &def.ColumnDef{},
			tag:  reflect.StructTag(`json:"abc"`),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.def.Tag = c.tag
			got := def.ParseColDef(def.WithModelTagKey(context.Background(), "db"), typ, c.tag)
			Expect(t, typ.String(), Equal(got.Type.String()))
			got.Type = nil
			Expect(t, got, Equal(c.def))
		})
	}
}
