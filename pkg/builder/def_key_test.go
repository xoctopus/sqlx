package builder_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func TestKeys(t *testing.T) {
	tab := builder.T("t_demo")
	var (
		cID        = builder.C("f_id", builder.WithColFieldName("ID"))
		cName      = builder.C("f_name", builder.WithColFieldName("Name"))
		cOrg       = builder.C("f_org", builder.WithColFieldName("Org"))
		cDeletedAt = builder.C("f_deleted_at", builder.WithColFieldName("DeletedAt"))
	)

	tab.(builder.ColsManager).AddCol(cID, cName, cOrg, cDeletedAt)
	tab.(builder.KeysManager).AddKey(
		builder.PK(builder.ColsOf(cID)),
		builder.UK(
			"ui_name",
			builder.ColsOf(tab.C("f_name")),
			builder.WithKeyMethod("BTREE"),
			builder.WithKeyColumnOptions(def.ResolveKeyColumnOptions("Name,NULL,FIRST")...),
		),
		builder.K(
			"i_org",
			builder.ColsOf(tab.C("f_org"), tab.C("DeletedAt")),
		),
	)

	t.Run("KeyColumnsDef", func(t *testing.T) {
		Expect(t, builder.KeyColumnsDefOf(tab.K("primary")), BeFragment("f_id"))
		Expect(t, builder.KeyColumnsDefOf(tab.K("ui_name")), BeFragment("f_name NULL FIRST"))
		Expect(t, builder.KeyColumnsDefOf(tab.K("i_org")), BeFragment("f_org,f_deleted_at"))
	})

	t.Run("K", func(t *testing.T) {
		Expect(t, tab.K("PRIMARY").Name(), Equal("primary"))
		Expect(t, tab.K("NON"), BeNil[builder.Key]())
		Expect(t, tab.K("PRIMARY").IsPrimary(), BeTrue())
		Expect(t, tab.K("ui_name").IsUnique(), BeTrue())
		Expect(t, tab.K("ui_name").(builder.KeyDef).Method(), Equal("BTREE"))
		Expect(
			t,
			tab.K("ui_name").(builder.KeyDef).ColumnOptions(),
			Equal([]builder.KeyColumnOption{{Name: "Name", Options: []string{"NULL", "FIRST"}}}),
		)
		Expect(t, tab.K("ui_name").IsNil(), BeFalse())
		Expect(t, tab.K("ui_name").String(), Equal("t_demo.ui_name"))
		Expect(t, frag.Func(tab.K("ui_name").Frag), BeFragment(tab.K("ui_name").Name()))
	})

}
