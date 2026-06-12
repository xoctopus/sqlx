package session

import (
	"context"
	"strings"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/builder"
)

// For retrieves Session by session name or table
func For(ctx context.Context, m any) (Session, bool) {
	schema := ""
	if x, ok := m.(builder.HasSchema); ok {
		schema = x.Schema()
	}

	switch x := m.(type) {
	case interface{ Unwrap() builder.Model }:
		return From(ctx, schema, x.Unwrap().TableName())
	case interface{ Unwrap() builder.Table }:
		return From(ctx, schema, x.Unwrap().TableName())
	case interface{ TableName() string }:
		return From(ctx, schema, x.TableName())
	case string:
		if before, after, ok := strings.Cut(x, "."); ok {
			schema = before
			x = after
		}
		return From(ctx, schema, x)
	default:
		return nil, false
	}
}

func MustFor(ctx context.Context, m any) Session {
	s, ok := For(ctx, m)
	must.BeTrueF(ok, "missing session for %T", m)
	return s
}

type tSession struct {
	schema string
	name   string
}

type SchemaModel interface {
	builder.HasSchema
	builder.HasTableName
}

// From retrieve Session from ctx by Session.Name
func From(ctx context.Context, schema, name string) (Session, bool) {
	s, ok := ctx.Value(tSession{schema, name}).(Session)
	return s, ok
}

func FromM(ctx context.Context, m SchemaModel) (Session, bool) {
	return From(ctx, m.Schema(), m.TableName())
}

func Must(ctx context.Context, schema, name string) Session {
	s, ok := From(ctx, schema, name)
	must.BeTrueF(ok, "missing session for: %s", name)
	return s
}

func MustM(ctx context.Context, m SchemaModel) Session {
	return Must(ctx, m.Schema(), m.TableName())
}

// With injects Session
func With(ctx context.Context, m builder.Model, s Session) context.Context {
	v := tSession{schema: s.Schema(), name: m.TableName()}
	return context.WithValue(ctx, v, s)
}

// Carry returns context carrier
func Carry(m builder.Model, s Session) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, m, s)
	}
}
