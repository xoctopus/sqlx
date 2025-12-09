package def_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
)

func TestParseKeyDef(t *testing.T) {
	cases := []struct {
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
					{Name: "Name", Options: []string{}},
				},
			},
		},
		{
			def: "primary ID ",
			opt: &def.KeyDefine{
				Kind: def.KEY_KIND__PRIMARY,
				Name: "primary",
				Options: []def.KeyColumnOption{
					{Name: "ID", Options: []string{}},
				},
			},
		},
		{
			def: " unique_index   idx_name   f_org_id,NULLS,FIRST;MemberID ",
			opt: &def.KeyDefine{
				Kind: def.KEY_KIND__UNIQUE_INDEX,
				Name: "idx_name",
				Options: []def.KeyColumnOption{
					{Name: "f_org_id", Options: []string{"NULLS", "FIRST"}},
					{Name: "MemberID", Options: []string{}},
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
	}
	for _, c := range cases {
		Expect(t, def.ParseKeyDef(c.def), Equal(c.opt))
	}

	Expect(t, cases[2].opt.OptionsNames(), Equal([]string{"f_org_id", "MemberID"}))
	Expect(t, cases[2].opt.OptionsStrings(), Equal([]string{
		"f_org_id,NULLS,FIRST",
		"MemberID",
	}))

	Expect(
		t,
		def.ResolveKeyColumnOptionsFromStrings(cases[2].opt.OptionsStrings()...),
		Equal(cases[2].opt.Options),
	)

	Expect(
		t,
		def.KeyColumnOptionByNames("f_id"),
		Equal([]def.KeyColumnOption{{Name: "f_id"}}),
	)
}
