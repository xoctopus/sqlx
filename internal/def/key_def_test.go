package def_test

import (
	"fmt"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
)

func ExampleParseKeyDef() {
	// // @model pk=ID
	kd := def.ParseKeyDef("pk", "ID")
	fmt.Printf("pk name=%s fields=%v\n", kd.Name, kd.OptionsNames())

	// // @model idx=ui_username,BTREE;Username
	kd = def.ParseKeyDef("idx", "ui_username,BTREE;Username")
	fmt.Printf("idx name=%s using=%s fields=%v\n", kd.Name, kd.Using, kd.OptionsNames())

	// // @model idx=i_status;Status,NULLS,FIRST
	kd = def.ParseKeyDef("idx", "i_status;Status,NULLS,FIRST")
	fmt.Printf("idx fields=%v\n", kd.OptionsStrings())

	// // @model uidx=ui_product_id;ProductID;DeletedAt
	kd = def.ParseKeyDef("uidx", "ui_product_id;ProductID;DeletedAt")
	fmt.Printf("uidx name=%s fields=%v\n", kd.Name, kd.OptionsNames())

	// Output:
	// pk name=primary fields=[ID]
	// idx name=ui_username using=BTREE fields=[Username]
	// idx fields=[Status,NULLS,FIRST]
	// uidx name=ui_product_id fields=[ProductID DeletedAt]
}

func TestParseKeyDef(t *testing.T) {
	cases := []struct {
		key string
		def string
		opt *def.KeyDefine
	}{
		{
			key: "idx",
			def: "idx_name,BTREE;Name",
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
			key: "pk",
			def: "ID",
			opt: &def.KeyDefine{
				Kind: def.KEY_KIND__PRIMARY,
				Name: "primary",
				Options: []def.KeyColumnOption{
					{Name: "ID", Options: []string{}},
				},
			},
		},
		{
			key: "uidx",
			def: "idx_name;f_org_id,NULLS,FIRST;MemberID ",
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
		{def: ""},
		{key: "pk", def: ""},
		{key: "idx", def: "idx_name"},
		{key: "uidx", def: "idx_name"},
		{key: "uidx", def: ";idx_name"},
		{key: "uidx", def: "idx_name;"},
		{key: "invalid", def: "idx_name;"},
	}
	for _, c := range cases {
		Expect(t, def.ParseKeyDef(c.key, c.def), Equal(c.opt))
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
