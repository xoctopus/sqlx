package builder_test

import (
	"context"
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/frag/testutil"
	"github.com/xoctopus/sqlx/testdata"
)

func TestTable(t *testing.T) {
	var (
		colID = builder.C(
			"f_id",
			builder.WithColFieldName("ID"),
			builder.WithColDefOf(uint64(0), ",autoinc"),
		)
		colName = builder.C(
			"f_name",
			builder.WithColFieldName("Name"),
			builder.WithColDefOf("", ",width=128,default=''"),
		)
	)
	tUser := builder.T(
		"t_user",
		colID, colName,
		builder.UK("u_idx_name", builder.ColsOf(colID, colName)),
		builder.PK(builder.ColsOf(colID)),
	)

	tUserRole := builder.T("t_user_role",
		builder.C(
			"f_id",
			builder.WithColFieldName("ID"),
			builder.WithColDefOf(uint64(0), ",autoinc"),
		),
		builder.C(
			"f_user_id",
			builder.WithColFieldName("UserID"),
			builder.WithColDefOf(uint64(0), ""),
		),
	)

	t.Run("Fragment", func(t *testing.T) {
		Expect(t, tUser.Fragment(""), BeNil[frag.Fragment]())

		Expect(
			t,
			tUser.Fragment("#.*"),
			testutil.BeFragment("t_user.*"),
		)
		Expect(
			t,
			tUser.Fragment("#ID = #ID + 1"),
			testutil.BeFragment("f_id = f_id + 1"),
		)
		Expect(
			t,
			tUser.Fragment("COUNT(#ID)"),
			testutil.BeFragment("COUNT(f_id)"),
		)
	})

	t.Run("HandleContext", func(t *testing.T) {
		f := builder.Select(nil).
			From(
				tUser,
				builder.Where(
					builder.AsCond(tUser.Fragment("#ID > 1")),
				),
				builder.Join(tUserRole).On(
					builder.AsCond(tUser.Fragment("#ID = ?", tUserRole.Fragment("#UserID"))),
				),
			)
		q := `SELECT * FROM t_user JOIN t_user_role ON t_user.f_id = t_user_role.f_user_id WHERE t_user.f_id > 1`
		Expect[frag.Fragment](t, f, testutil.BeFragment(q))
	})

	t.Run("WithTableName", func(t *testing.T) {
		tUserOf := tUser.(builder.WithTableName).WithTableName("t_user_2")
		Expect(t, tUserOf.TableName(), Equal("t_user_2"))
		Expect(t, tUserOf.String(), Equal("t_user_2"))
	})

	t.Run("Picking", func(t *testing.T) {
		col := tUser.C("ID")
		Expect(t, col.String(), Equal("t_user.f_id"))

		Expect(t, slices.Collect(tUser.Cols()), HaveLen[[]builder.Col](2))

		uk := tUser.K("u_idx_name")
		Expect(t, uk.IsPrimary(), BeFalse())
		Expect(t, uk.IsUnique(), BeTrue())
		Expect(t, uk.(builder.WithTable).T(), Equal(tUser))
		pk := tUser.K("PRIMARY")
		Expect(t, pk.IsPrimary(), BeTrue())
		Expect(t, pk.IsUnique(), BeTrue())
		Expect(t, pk.(builder.WithTable).T(), Equal(tUser))

		Expect(t, slices.Collect(tUser.Keys()), HaveLen[[]builder.Key](2))

		Expect(t, tUser.Pick("f_id").C("f_id"), Equal(tUser.C("f_id")))
	})

	t.Run("Namespaces", func(t *testing.T) {
		tab := builder.T("t").(builder.WithSchema).WithSchema("schema").(builder.WithDatabase).WithDatabase("database")
		Expect(t, tab.(builder.HasSchema).Schema(), Equal("schema"))
		Expect(t, tab.(builder.HasDatabase).Database(), Equal("database"))
	})
}

type Role struct {
}

type WithAttrs struct {
	ID int64
}

func (WithAttrs) ColumnComment() map[string]string {
	return map[string]string{
		"ID": "autoinc pk",
	}
}

func (WithAttrs) ColumnDesc() map[string][]string {
	return map[string][]string{
		"ID": {"desc line1", "desc line2"},
	}
}

func (WithAttrs) ColumnRel() map[string][]string {
	return map[string][]string{
		"ID": {"t_a.f_a_id", "TB.FieldID"},
	}
}

func TestCatalog(t *testing.T) {
	ctx := context.Background()
	catalog := builder.CatalogFrom(
		ctx,
		&testdata.User{},
		&testdata.Org{},
		builder.TFrom(ctx, &Role{}),
	)

	Expect(t, builder.NewCatalog().Len(), Equal(0))

	t.Run("Pick", func(t *testing.T) {
		Expect(
			t,
			slices.Collect(builder.TableNames(catalog)),
			EquivalentSlice([]string{"Role", "users", "t_org"}),
		)
		Expect(t, catalog.T("users").TableName(), Equal("users"))
		Expect(t, catalog.T("exclude"), BeNil[builder.Table]())
	})

	t.Run("Replace", func(t *testing.T) {
		catalog.Add(builder.TFrom(ctx, &Role{}))
		Expect(t, slices.Collect(catalog.Tables()), HaveLen[[]builder.Table](catalog.Len()))
	})

	t.Run("WithAttrs", func(t *testing.T) {
		m := &WithAttrs{}
		tab := builder.TFrom(ctx, m)
		d := tab.C("ID").(builder.ColDef).Def()
		Expect(t, d.Relation, Equal(m.ColumnRel()["ID"]))
		Expect(t, d.Comment, Equal(m.ColumnComment()["ID"]))
		Expect(t, d.Desc, Equal(m.ColumnDesc()["ID"]))
	})
}
