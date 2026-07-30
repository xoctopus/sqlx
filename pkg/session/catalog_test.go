package session

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
)

func TestRegister(t *testing.T) {
	t.Cleanup(func() { catalogs.Clear() })

	t.Run("OK", func(t *testing.T) {
		catalogs.Clear()

		a := builder.NewCatalog()
		a.Add(builder.T("t_a"))
		a.Add(builder.T("t_b"))

		b := builder.NewCatalog()
		b.Add(builder.T("t_user").(builder.WithSchema).WithSchema("app"))

		Register(a, b)
		_, ok := catalogs.Load(".t_a")
		Expect(t, ok, BeTrue())
		_, ok = catalogs.Load(".t_b")
		Expect(t, ok, BeTrue())
		_, ok = catalogs.Load("app.t_user")
		Expect(t, ok, BeTrue())
	})

	t.Run("SameNameDifferentSchema", func(t *testing.T) {
		catalogs.Clear()

		c1 := builder.NewCatalog()
		c1.Add(builder.T("t_user").(builder.WithSchema).WithSchema("s1"))
		c2 := builder.NewCatalog()
		c2.Add(builder.T("t_user").(builder.WithSchema).WithSchema("s2"))

		Register(c1)
		Register(c2)
		_, ok := catalogs.Load("s1.t_user")
		Expect(t, ok, BeTrue())
		_, ok = catalogs.Load("s2.t_user")
		Expect(t, ok, BeTrue())
	})

	t.Run("DuplicatePanics", func(t *testing.T) {
		catalogs.Clear()

		cat := builder.NewCatalog()
		cat.Add(builder.T("t_dup"))
		Register(cat)

		ExpectPanic[error](t, func() {
			Register(cat)
		}, ErrorContains("already registered"))
	})

	t.Run("DuplicateAcrossCatalogs", func(t *testing.T) {
		catalogs.Clear()

		c1 := builder.NewCatalog()
		c1.Add(builder.T("t_shared"))
		c2 := builder.NewCatalog()
		c2.Add(builder.T("t_shared"))

		Register(c1)
		ExpectPanic[error](t, func() {
			Register(c2)
		}, ErrorContains("already registered"))
	})
}
