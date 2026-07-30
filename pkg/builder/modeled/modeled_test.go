package modeled_test

import (
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/builder/modeled"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/testdata"
)

func TestModeled(t *testing.T) {
	m := modeled.M[testdata.User]()

	t.Run("Table", func(t *testing.T) {
		Expect(t, m.TableName(), Equal("t_user"))
		mm := m.Model()
		Expect(t, mm != nil, BeTrue())

		raw := builder.TFrom(&testdata.User{})
		casted := modeled.CastT[testdata.User](raw)
		Expect(t, casted.TableName(), Equal(raw.TableName()))
		Expect(t, casted.C("ID").Name(), Equal(raw.C("ID").Name()))
	})

	t.Run("MCols", func(t *testing.T) {
		cols := slices.Collect(m.Cols())
		mcols := slices.Collect(m.MCols())
		Expect(t, mcols, HaveLen[[]modeled.Col[testdata.User]](len(cols)))

		for i, c := range mcols {
			Expect(t, c.Name(), Equal(cols[i].Name()))
		}
	})

	t.Run("MKeys", func(t *testing.T) {
		keys := slices.Collect(m.Keys())
		mkeys := slices.Collect(m.MKeys())
		Expect(t, mkeys, HaveLen[[]modeled.Key[testdata.User]](len(keys)))

		for i, k := range mkeys {
			Expect(t, k.Name(), Equal(keys[i].Name()))
		}
	})

	t.Run("MK", func(t *testing.T) {
		pk := m.MK("primary")
		Expect(t, pk.IsPrimary(), BeTrue())
		Expect(t, pk.IsUnique(), BeTrue())

		uk := m.MK("ui_name")
		Expect(t, uk.IsPrimary(), BeFalse())
		Expect(t, uk.IsUnique(), BeTrue())

		pkCols := slices.Collect(pk.MCols())
		Expect(t, pkCols, HaveLen[[]modeled.Col[testdata.User]](1))
		Expect(t, pkCols[0].Name(), Equal(m.C("ID").Name()))
	})

	t.Run("CastC", func(t *testing.T) {
		raw := m.C("Name")
		c := modeled.CastC[testdata.User](raw)
		Expect(t, c.Name(), Equal(raw.Name()))

		computed := c.ComputedBy(frag.Lit("UPPER(f_name)"))
		Expect(t, computed.Name(), Equal(raw.Name()))
	})

	t.Run("CT", func(t *testing.T) {
		raw := m.C("ID")
		c := modeled.CT[testdata.User, uint64](raw)
		Expect(t, c.Name(), Equal(raw.Name()))

		computed := c.ComputedBy(builder.Distinct(raw))
		Expect(t, computed.Name(), Equal(raw.Name()))

		typed := c.TypedComputedBy(builder.Distinct(raw))
		Expect(t, typed.Name(), Equal(raw.Name()))
	})

	t.Run("CastK", func(t *testing.T) {
		raw := m.K("i_age")
		k := modeled.CastK[testdata.User](raw)
		Expect(t, k.Name(), Equal(raw.Name()))
		Expect(t, k.IsUnique(), BeFalse())

		cols := slices.Collect(k.MCols())
		Expect(t, len(cols), Equal(1))
		Expect(t, cols[0].Name(), Equal(m.C("Age").Name()))
	})
}
