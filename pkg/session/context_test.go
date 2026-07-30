package session_test

import (
	"context"
	"database/sql"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/session"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
)

type stubSession struct {
	name string
}

func (s *stubSession) Name() string { return s.name }
func (s *stubSession) T(any) builder.Table {
	return nil
}
func (s *stubSession) Adaptor(...session.AdaptorOptionApplier) adaptor.Adaptor {
	return nil
}
func (s *stubSession) Tx(context.Context, func(context.Context) error) error {
	return nil
}
func (s *stubSession) Exec(context.Context, frag.Fragment) (sql.Result, error) {
	return nil, nil
}
func (s *stubSession) Query(context.Context, frag.Fragment) (*sql.Rows, error) {
	return nil, nil
}

type demoModel struct{}

func (demoModel) TableName() string { return "t_demo" }

type demoModelMissing struct{}

func (demoModelMissing) TableName() string { return "t_missing" }

type schemaModel struct {
	schema string
}

func (m schemaModel) TableName() string { return "t_demo" }
func (m schemaModel) Schema() string    { return m.schema }

func TestContext(t *testing.T) {
	s := &stubSession{name: "main"}
	m := demoModel{}
	sm := schemaModel{schema: "app"}

	t.Run("WithAndFrom", func(t *testing.T) {
		ctx := session.With(context.Background(), s)

		got, ok := session.From(ctx, "main")
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		_, ok = session.From(ctx, "missing")
		Expect(t, ok, BeFalse())

		Expect(t, session.Must(ctx, "main"), Be[session.Session](s))
		ExpectPanic[error](t, func() {
			session.Must(ctx, "missing")
		}, ErrorContains("missing session"))
	})

	t.Run("WithModel", func(t *testing.T) {
		ctx := session.WithModel(context.Background(), m, s)

		got, ok := session.From(ctx, "t_demo")
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		got, ok = session.FromByModel(ctx, m)
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		Expect(t, session.MustByModel(ctx, m), Be[session.Session](s))
		ExpectPanic[error](t, func() {
			session.MustByModel(ctx, demoModelMissing{})
		}, ErrorContains("missing session"))
	})

	t.Run("WithSchemaModel", func(t *testing.T) {
		ctx := session.WithSchemaModel(context.Background(), sm, s)

		_, ok := session.From(ctx, "t_demo")
		Expect(t, ok, BeFalse())

		got, ok := session.From(ctx, "app.t_demo")
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		got, ok = session.FromByModel(ctx, sm)
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		// falls back to TableName when schema key misses
		ctx2 := session.WithModel(context.Background(), sm, s)
		got, ok = session.FromByModel(ctx2, sm)
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))
	})

	t.Run("For", func(t *testing.T) {
		ctx := session.With(context.Background(), s)
		ctx = session.WithModel(ctx, m, s)

		got, ok := session.For(ctx, "main")
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		got, ok = session.For(ctx, m)
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		_, ok = session.For(ctx, 123)
		Expect(t, ok, BeFalse())

		Expect(t, session.MustFor(ctx, "main"), Be[session.Session](s))
		Expect(t, session.MustFor(ctx, m), Be[session.Session](s))
		ExpectPanic[error](t, func() {
			session.MustFor(ctx, 123)
		}, ErrorContains("missing session"))
	})

	t.Run("Carry", func(t *testing.T) {
		ctx := session.CarryModel(m, s)(context.Background())
		got, ok := session.FromByModel(ctx, m)
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))

		ctx = session.CarrySchemaModel(sm, s)(context.Background())
		got, ok = session.FromByModel(ctx, sm)
		Expect(t, ok, BeTrue())
		Expect(t, got, Be[session.Session](s))
	})
}
