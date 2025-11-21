package builder_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func TestColumns(t *testing.T) {
	cs := builder.Columns()

	Expect(t, cs.Len(), Equal(0))

	cs.(builder.ColsManager).AddCol(
		builder.C(
			"f_id",
			builder.WithColFieldName("ID"),
			builder.WithColDefOf(100, `,autoinc`),
		),
		builder.C(
			"f_name",
			builder.WithColFieldName("Name"),
			builder.WithColDefOf("saito", ``),
		),
	)
	cs = cs.Of(builder.T("t_table"))

	t.Run("Add", func(t *testing.T) {
		t.Run("GetByFieldName", func(t *testing.T) {
			c := cs.C("ID")
			Expect(t, c.Name(), Equal("f_id"))
			Expect(t, c.FieldName(), Equal("ID"))

			picked := builder.PickCols(cs, "ID", "Name")
			Expect(t, picked.Len(), Equal(2))

			ExpectPanic[error](t, func() {
				builder.PickCols(cs, "unknown").Len()
			}, ErrorContains("unknown column"))
		})
		t.Run("GetByColumnName", func(t *testing.T) {
			c := cs.C("f_id")
			Expect(t, c.Name(), Equal("f_id"))
			Expect(t, c.FieldName(), Equal("ID"))

			picked := builder.PickCols(cs, "f_id", "f_name")
			Expect(t, picked.Len(), Equal(2))
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
		c1 := builder.CastC[int](cs.C("ID"))
		c2 := builder.CastC[int](builder.C("f_other_id"))
		c3 := builder.CastC[string](builder.C("f_name"))

		Expect(t, c1.AsCond(nil), BeNil[frag.Fragment]())
		Expect(t, c1.AssignBy(), BeNil[builder.Assignment]())

		t.Run("EqNeq", func(t *testing.T) {
			q, args := frag.Collect(bgc, c1.AsCond(builder.Eq(1)))
			Expect(t, q, Equal("f_id = ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.Neq(1)))
			Expect(t, q, Equal("f_id <> ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.EqCol(c2)))
			Expect(t, q, Equal("f_id = f_other_id"))
			Expect(t, args, HaveLen[[]any](0))

			q, args = frag.Collect(bgc, c1.AsCond(builder.NeqCol(c2)))
			Expect(t, q, Equal("f_id <> f_other_id"))
			Expect(t, args, HaveLen[[]any](0))
		})
		t.Run("In", func(t *testing.T) {
			q, args := frag.Collect(bgc, c1.AsCond(builder.In(1, 2, 3)))
			Expect(t, q, Equal("f_id IN (?,?,?)"))
			Expect(t, args, Equal([]any{1, 2, 3}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.NotIn(1, 2, 3)))
			Expect(t, q, Equal("f_id NOT IN (?,?,?)"))
			Expect(t, args, Equal([]any{1, 2, 3}))

			q, args = frag.Collect(bgc, c1.AsCond(builder.In[int]()))
			Expect(t, q, Equal(""))
			Expect(t, args, HaveLen[[]any](0))

			q, args = frag.Collect(bgc, c1.AsCond(builder.NotIn[int]()))
			Expect(t, q, Equal(""))
			Expect(t, args, HaveLen[[]any](0))
		})
		t.Run("Null", func(t *testing.T) {
			q, args := frag.Collect(bgc, c1.AsCond(builder.IsNull[int]()))
			Expect(t, q, Equal("f_id IS NULL"))
			Expect(t, args, HaveLen[[]any](0))

			q, args = frag.Collect(bgc, c1.AsCond(builder.IsNotNull[int]()))
			Expect(t, q, Equal("f_id IS NOT NULL"))
			Expect(t, args, HaveLen[[]any](0))
		})
		t.Run("Like", func(t *testing.T) {
			q, args := frag.Collect(bgc, c3.AsCond(builder.Like[string]("1")))
			Expect(t, q, Equal("f_name LIKE ?"))
			Expect(t, args, Equal([]any{"%1%"}))

			q, args = frag.Collect(bgc, c3.AsCond(builder.LLike[string]("1")))
			Expect(t, q, Equal("f_name LIKE ?"))
			Expect(t, args, Equal([]any{"%1"}))

			q, args = frag.Collect(bgc, c3.AsCond(builder.RLike[string]("1")))
			Expect(t, q, Equal("f_name LIKE ?"))
			Expect(t, args, Equal([]any{"1%"}))

			q, args = frag.Collect(bgc, c3.AsCond(builder.NotLike[string]("1")))
			Expect(t, q, Equal("f_name NOT LIKE ?"))
			Expect(t, args, Equal([]any{"%1%"}))
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
			q, args := frag.Collect(bgc, c1.AssignBy(builder.AsValue(c2)))
			Expect(t, q, Equal("f_id = f_other_id"))
			Expect(t, args, HaveLen[[]any](0))

			q, args = frag.Collect(bgc, builder.ColumnsAndValues(
				builder.Columns("f_id", "f_name"),
				builder.AsValue(c1), builder.AsValue(c3),
			))
			Expect(t, q, Equal("(f_id,f_name) VALUES (?,?)"))
			Expect(t, args, HaveLen[[]any](2))

			q, args = frag.Collect(bgc, c1.AssignBy(builder.Value(1)))
			Expect(t, q, Equal("f_id = ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AssignBy(builder.Inc(1)))
			Expect(t, q, Equal("f_id = f_id + ?"))
			Expect(t, args, Equal([]any{1}))

			q, args = frag.Collect(bgc, c1.AssignBy(builder.Dec(1)))
			Expect(t, q, Equal("f_id = f_id - ?"))
			Expect(t, args, Equal([]any{1}))
		})
		t.Run("Fragment", func(t *testing.T) {
			q, args := frag.Collect(bgc, c1.Fragment("# = ?", 1))
			Expect(t, q, Equal("f_id = ?"))
			Expect(t, args, Equal([]any{1}))
		})

		t.Run("Frag", func(t *testing.T) {
			t.Run("InProject", func(t *testing.T) {
				cc := builder.CastC[int](
					builder.C("distinct_ids"),
					builder.WithColComputed(builder.Distinct(c1)),
				)

				q, args := frag.Collect(
					builder.WithToggles(bgc, builder.TOGGLE__IN_PROJECT),
					frag.Func(cc.Frag),
				)
				Expect(t, q, Equal("DISTINCT(f_id) AS distinct_ids"))
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
	})
}

func BenchmarkCols(b *testing.B) {
	cols := builder.Columns()

	(cols).(builder.ColsManager).AddCol(
		builder.C("f_id", builder.WithColFieldName("ID"), builder.WithColDefOf(1, `,autoinc`)),
		builder.C("f_name", builder.WithColFieldName("Name"), builder.WithColDefOf(1, ``)),
		builder.C("f_f1", builder.WithColFieldName("F1"), builder.WithColDefOf(1, ``)),
		builder.C("f_f2", builder.WithColFieldName("F2"), builder.WithColDefOf(1, ``)),
		builder.C("f_f3", builder.WithColFieldName("F3"), builder.WithColDefOf(1, ``)),
		builder.C("f_f4", builder.WithColFieldName("F4"), builder.WithColDefOf(1, ``)),
		builder.C("f_f5", builder.WithColFieldName("F5"), builder.WithColDefOf(1, ``)),
		builder.C("f_f6", builder.WithColFieldName("F6"), builder.WithColDefOf(1, ``)),
		builder.C("f_f7", builder.WithColFieldName("F7"), builder.WithColDefOf(1, ``)),
		builder.C("f_f8", builder.WithColFieldName("F8"), builder.WithColDefOf(1, ``)),
		builder.C("f_f9", builder.WithColFieldName("F9"), builder.WithColDefOf(1, ``)),
	)

	b.Run("Single", func(b *testing.B) {
		for b.Loop() {
			_ = cols.C("F3")
		}
	})

	b.Run("Multi", func(b *testing.B) {
		for b.Loop() {
			_ = builder.PickCols(cols, "ID", "Name")
		}
	})
}
