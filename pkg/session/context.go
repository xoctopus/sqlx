package session

import (
	"context"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/builder"
)

type ModelWrapper interface {
	Unwrap() builder.Table
}

// For retrieves Session by session name or table
func For(ctx context.Context, m any) (Session, bool) {
	k := ""
	if x, ok := m.(builder.HasSession); ok {
		k = x.Session()
	}
	if len(k) == 0 {
		k, _ = KeyFrom(ctx)
	}
	if len(k) == 0 {
		return nil, false
	}

	switch x := m.(type) {
	case interface{ Unwrap() builder.Model }:
		return From(ctx, k+"."+x.Unwrap().TableName())
	case interface{ Unwrap() builder.Table }:
		return From(ctx, k+"."+x.Unwrap().TableName())
	case interface{ TableName() string }:
		return From(ctx, k+"."+x.TableName())
	case string:
		return From(ctx, k+"."+x)
	default:
		return nil, false
	}
}

func MustFor(ctx context.Context, m any) Session {
	s, ok := For(ctx, m)
	must.BeTrueF(ok, "missing session for %T", m)
	return s
}

type (
	tSessionKey struct{}
	tSession    struct{ name string }
)

var (
	KeyFrom  = contextx.From[tSessionKey, string]
	KeyMust  = contextx.Must[tSessionKey, string]
	KeyWith  = contextx.With[tSessionKey, string]
	KeyCarry = contextx.Carry[tSessionKey, string]
)

// From retrieve Session from ctx by Session.Name
func From(ctx context.Context, name string) (Session, bool) {
	s, ok := ctx.Value(tSession{name}).(Session)
	return s, ok
}

func Must(ctx context.Context, name string) Session {
	s, ok := From(ctx, name)
	must.BeTrueF(ok, "missing session for: %s", name)
	return s
}

// With injects Session
func With(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, tSession{name: session.Name()}, session)
}

// Carry returns context carrier
func Carry(session Session) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, session)
	}
}
