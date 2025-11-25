package def_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
)

func TestParseKeyDef(t *testing.T) {
	for _, c := range []struct {
		def string
		opt *def.KeyDefine
	}{
		{
			def: "idx idx_name,BTREE Name",
			opt: &def.KeyDefine{
				Kind:  def.KEY_KIND__INDEX,
				Name:  "idx_name",
				Using: "BTREE",
				Options: []def.KeyColumnOption{
					{FieldName: "Name", Options: []string{}},
				},
			},
		},
		{
			def: "primary ID ",
			opt: &def.KeyDefine{
				Kind: def.KEY_KIND__PRIMARY,
				Options: []def.KeyColumnOption{
					{FieldName: "ID", Options: []string{}},
				},
			},
		},
		{
			def: " unique_index   idx_name   OrgID,NULLS,FIRST;MemberID ",
			opt: &def.KeyDefine{
				Kind: def.KEY_KIND__UNIQUE_INDEX,
				Name: "idx_name",
				Options: []def.KeyColumnOption{
					{FieldName: "OrgID", Options: []string{"NULLS", "FIRST"}},
					{FieldName: "MemberID", Options: []string{}},
				},
			},
		},
		// invalid
		{def: "index"},
		{def: "u_idx"},
		{def: "pk"},
		{def: "index ,using ID"},
		{def: "index idx_name,using ;"},
		{def: "invalid"},
	} {
		Expect(t, def.ParseKeyDef(c.def), Equal(c.opt))
	}
}
