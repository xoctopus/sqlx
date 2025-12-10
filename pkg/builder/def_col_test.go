package builder_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/builder/modeled"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/frag/testutil"
)

type Fragment = frag.Fragment

var (
	BeFragment         = testutil.BeFragment
	BeFragmentForQuery = testutil.BeFragmentForQuery
)

func TestColumns(t *testing.T) {
	cs := builder.Columns()

	Expect(t, cs.Len(), Equal(0))
	Expect(t, cs.C(""), BeNil[builder.Col]())

	cs.(builder.ColsManager).AddCol(
		builder.C(
			"f_id",
			builder.WithColFieldName("ID"),
			builder.WithColDefOf(context.Background(), 100, `,autoinc`),
		),
		builder.C(
			"f_name",
			builder.WithColFieldName("Name"),
			builder.WithColDefOf(context.Background(), "", ``),
		),
	)
	cs = cs.Of(builder.T("t_table"))

	cs2 := builder.ColsIterOf(cs.Cols())
	Expect(t, cs2.Len(), Equal(cs.Len()))

	t.Run("Add", func(t *testing.T) {
		t.Run("GetByFieldName", func(t *testing.T) {
			c := cs.C("ID")
			Expect(t, c.Name(), Equal("f_id"))
			Expect(t, c.FieldName(), Equal("ID"))

			sub := cs.Pick("ID", "Name")
			Expect(t, sub.Len(), Equal(2))

			ExpectPanic[error](t, func() { cs.Pick("unknown") }, ErrorContains("unknown column"))
		})
		t.Run("GetByColumnName", func(t *testing.T) {
			c := cs.C("f_id")
			Expect(t, c.Name(), Equal("f_id"))
			Expect(t, c.FieldName(), Equal("ID"))

			sub := cs.Pick("f_id", "f_name")
			Expect(t, sub.Len(), Equal(2))
		})
	})
	t.Run("Of", func(t *testing.T) {
		c := builder.C("f_id")
		Expect(t, c.String(), Equal(c.Name()))
		c = c.Of(builder.T("t_table"))
		Expect(t, c.String(), Equal("t_table.f_id"))
	})
	t.Run("Compute", func(t *testing.T) {
		bgc := context.Background()
		c1 := builder.CC[int](cs.C("ID").Of(builder.T("t1")))
		c2 := builder.CC[int](builder.C("f_other_id").Of(builder.T("t1")))
		c3 := builder.CC[string](builder.C("f_name").Of(builder.T("t1")))

		Expect(t, c1.AsCond(nil), BeNil[Fragment]())
		Expect(t, c1.AssignBy(), BeNil[builder.Assignment]())

		t.Run("EqNeq", func(t *testing.T) {
			Expect(t, c1.AsCond(builder.Eq(1)), BeFragment("f_id = ?", 1))
			Expect(t, c1.AsCond(builder.Neq(1)), BeFragment("f_id <> ?", 1))
			Expect(t, c1.AsCond(builder.EqCol(c2)), BeFragment("f_id = f_other_id"))
			Expect(t, c1.AsCond(builder.NeqCol(c2)), BeFragment("f_id <> f_other_id"))
		})
		t.Run("In", func(t *testing.T) {
			Expect(t, c1.AsCond(builder.In(1, 2, 3)), BeFragment("f_id IN (?,?,?)", 1, 2, 3))
			Expect(t, c1.AsCond(builder.NotIn(1, 2, 3)), BeFragment("f_id NOT IN (?,?,?)", 1, 2, 3))
			Expect(t, c1.AsCond(builder.In[int]()), BeNil[Fragment]())
			Expect(t, c1.AsCond(builder.NotIn[int]()), BeNil[Fragment]())
		})
		t.Run("Null", func(t *testing.T) {
			Expect(t, c1.AsCond(builder.IsNull[int]()), BeFragment("f_id IS NULL"))
			Expect(t, c1.AsCond(builder.IsNotNull[int]()), BeFragment("f_id IS NOT NULL"))
		})
		t.Run("Like", func(t *testing.T) {
			Expect(t, c3.AsCond(builder.Like[string]("1")), BeFragment("f_name LIKE ?", "%1%"))
			Expect(t, c3.AsCond(builder.LLike[string]("1")), BeFragment("f_name LIKE ?", "%1"))
			Expect(t, c3.AsCond(builder.RLike[string]("1")), BeFragment("f_name LIKE ?", "1%"))
			Expect(t, c3.AsCond(builder.NotLike[string]("1")), BeFragment("f_name NOT LIKE ?", "%1%"))
		})
		t.Run("Between", func(t *testing.T) {
			q, args := frag.Collect(bgc, c1.AsCond(builder.Between(1, 2)))
			Expect(t, q, Equal("f_id BETWEEN ? AND ?"))
			Expect(t, args, Equal([]any{1, 2}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.NotBetween(1, 2)))
			Expect(t, q, Equal("f_id NOT BETWEEN ? AND ?"))
			Expect(t, args, Equal([]any{1, 2}))
		})
		t.Run("Compare", func(t *testing.T) {
			q, args := frag.Collect(bgc, c1.AsCond(builder.Gt(1)))
			Expect(t, q, Equal("f_id > ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.Gte(1)))
			Expect(t, q, Equal("f_id >= ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.Lt(1)))
			Expect(t, q, Equal("f_id < ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.Lte(1)))
			Expect(t, q, Equal("f_id <= ?"))
			Expect(t, args, Equal([]any{1}))
		})
		t.Run("Assign", func(t *testing.T) {
			Expect[Fragment](t, c1.AssignBy(builder.AsValue(c2)), BeFragment("f_id = f_other_id"))

			c1_ := c1.Of(builder.T("t2"))
			c3_ := c3.Of(builder.T("t2"))
			f := builder.ColumnsAndValues(builder.ColsOf(c1, c3), c1_, c3_)

			q, args := frag.Collect(
				builder.WithToggles(bgc, builder.TOGGLE__MULTI_TABLE),
				f,
			)
			Expect(t, q, Equal("(f_id,f_name) VALUES (t2.f_id,t2.f_name)"))
			Expect(t, len(args), Equal(0))

			Expect[Fragment](t, c1.AssignBy(builder.Value(1)), BeFragment("f_id = ?", 1))
			Expect[Fragment](t, c1.AssignBy(builder.Inc(1)), BeFragment("f_id = f_id + ?", 1))
			Expect[Fragment](t, c1.AssignBy(builder.Dec(1)), BeFragment("f_id = f_id - ?", 1))
		})
		t.Run("Fragment", func(t *testing.T) {
			Expect(t, c1.Fragment("# = ?", 1), BeFragment("f_id = ?", 1))
		})
		t.Run("Frag", func(t *testing.T) {
			t.Run("InProject", func(t *testing.T) {
				cc := builder.CC[int](
					builder.C("distinct_id"),
					builder.WithColComputed(builder.Distinct(c1)),
				)

				q, args := frag.Collect(
					builder.WithToggles(bgc, builder.TOGGLE__IN_PROJECT),
					frag.Func(cc.Frag),
				)
				Expect(t, q, Equal("DISTINCT(f_id) AS distinct_id"))
				Expect(t, args, HaveLen[[]any](0))
			})
			t.Run("MultiTable", func(t *testing.T) {
				cc := builder.CT[int]("f_id", builder.WithColFieldName("ID")).
					Of(builder.T("t_table"))

				q, args := frag.Collect(
					builder.WithToggles(bgc, builder.TOGGLE__MULTI_TABLE),
					frag.Func(cc.Frag),
				)
				Expect(t, q, Equal("t_table.f_id"))
				Expect(t, args, HaveLen[[]any](0))
				t.Run("AutoAlias", func(t *testing.T) {
					q, args = frag.Collect(
						builder.WithToggles(
							bgc,
							builder.TOGGLE__MULTI_TABLE,
							builder.TOGGLE__AUTO_ALIAS,
						),
						frag.Func(cc.Frag),
					)
					Expect(t, q, Equal("t_table.f_id AS t_table__f_id"))
					Expect(t, args, HaveLen[[]any](0))
				})
				Expect(t, cc.(builder.WithTable).T().TableName(), Equal("t_table"))
			})
		})

		t.Run("Modeled", func(t *testing.T) {
			cc := modeled.CT[builder.Model, int](
				builder.C(
					"distinct_id",
					builder.WithColComputed(builder.Avg(builder.C("sum"))),
					builder.WithColFieldName("SumAverage"),
				).Of(builder.T("t_demo")),
			)
			Expect(t, builder.GetColComputed(cc), BeFragment("AVG(sum)"))
			Expect(t, builder.GetColDef(cc), Equal(builder.ColumnDef{}))
			Expect(t, builder.GetColTable(cc).TableName(), Equal("t_demo"))
		})
	})
}

func BenchmarkCols(b *testing.B) {
	cols := builder.Columns()

	(cols).(builder.ColsManager).AddCol(
		builder.C("f_id", builder.WithColFieldName("ID"), builder.WithColDefOf(context.Background(), 1, `,autoinc`)),
		builder.C("f_name", builder.WithColFieldName("Name"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f1", builder.WithColFieldName("F1"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f2", builder.WithColFieldName("F2"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f3", builder.WithColFieldName("F3"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f4", builder.WithColFieldName("F4"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f5", builder.WithColFieldName("F5"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f6", builder.WithColFieldName("F6"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f7", builder.WithColFieldName("F7"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f8", builder.WithColFieldName("F8"), builder.WithColDefOf(context.Background(), 1, ``)),
		builder.C("f_f9", builder.WithColFieldName("F9"), builder.WithColDefOf(context.Background(), 1, ``)),
	)

	b.Run("Single", func(b *testing.B) {
		for b.Loop() {
			_ = cols.C("F3")
		}
	})

	b.Run("Multi", func(b *testing.B) {
		for b.Loop() {
			_ = cols.Pick("ID", "Name")
		}
	})
}
