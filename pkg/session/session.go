package session

import (
	"context"
	"database/sql"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
	_ "github.com/xoctopus/sqlx/pkg/sql/adaptor/mysql"
	_ "github.com/xoctopus/sqlx/pkg/sql/adaptor/postgres"
	_ "github.com/xoctopus/sqlx/pkg/sql/adaptor/sqlite"
)

// Session defines logic session interface
type Session interface {
	// Name logic session name
	Name() string
	// T picks table from session
	T(any) builder.Table
	// Adaptor returns session adaptor
	Adaptor(...AdaptorOptionApplier) adaptor.Adaptor

	Tx(context.Context, func(context.Context) error) error
	Exec(context.Context, frag.Fragment) (sql.Result, error)
	Query(context.Context, frag.Fragment) (*sql.Rows, error)
}

func New(a adaptor.Adaptor, name string) Session {
	return &session{
		name: name,
		a:    a,
	}
}

func NewReadonly(rw adaptor.Adaptor, ro adaptor.Adaptor, name string) Session {
	return &session{
		name: name,
		a:    rw,
		ro:   ro,
	}
}

type session struct {
	name string

	a  adaptor.Adaptor
	ro adaptor.Adaptor
}

func (s *session) Name() string {
	return s.name
}

func (s *session) T(m any) builder.Table {
	switch x := m.(type) {
	case builder.WithTable:
		return x.T()
	case builder.Table:
		return x
	default:
		return builder.TFrom(m)
	}
}

func (s *session) Tx(ctx context.Context, exec func(context.Context) error) error {
	return s.a.Tx(ctx, exec)
}

func (s *session) Exec(ctx context.Context, f frag.Fragment) (sql.Result, error) {
	return s.a.Exec(ctx, f)
}

func (s *session) Query(ctx context.Context, f frag.Fragment) (*sql.Rows, error) {
	return s.a.Query(ctx, f)
}

func (s *session) Adaptor(appliers ...AdaptorOptionApplier) adaptor.Adaptor {
	opt := &AdaptorOption{}
	for _, apply := range appliers {
		apply(opt)
	}

	if opt.ReadOnly {
		return s.ro
	}
	return s.a
}

type AdaptorOption struct {
	ReadOnly bool
}

type AdaptorOptionApplier func(*AdaptorOption)

func ReadOnly() AdaptorOptionApplier {
	return func(o *AdaptorOption) {
		o.ReadOnly = true
	}
}
