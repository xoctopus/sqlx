package session

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
)

// Session is a logical isolation unit and operational handle for database
// adapters.
// eg:
//
//	a specific MySQL database
//	a specific search_path in same PostgreSQL database
//	a particular SQLite database file.
type Session interface {
	// Schema logically isolation
	Schema() string
	// T picks table from session
	T(any) builder.Table
	// Tx exec query
	Tx(context.Context, func(context.Context) error) error
	// Adaptor returns session adaptor
	Adaptor(...AdaptorOptionApplier) adaptor.Adaptor
}

func New(a adaptor.Adaptor, schema string) Session {
	return &session{
		schema: a.Schema(),
		a:      a,
	}
}

func NewReadonly(rw adaptor.Adaptor, ro adaptor.Adaptor, name string) Session {
	return &session{
		schema: ro.Schema(),
		a:      rw,
		ro:     ro,
	}
}

type session struct {
	schema string
	// curr    string
	// escaped atomic.Bool

	a  adaptor.Adaptor
	ro adaptor.Adaptor
}

func (s *session) Schema() string {
	return s.schema
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
