package session_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/session"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
)

func TestInTx(t *testing.T) {
	Expect(t, session.InTx(context.Background()), BeFalse())

	db, mock, err := sqlmock.New()
	Expect(t, err, Succeed())
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectBegin()
	tx, err := db.Begin()
	Expect(t, err, Succeed())

	ctx := adaptor.WithExecutor(context.Background(), tx)
	Expect(t, session.InTx(ctx), BeTrue())

	ctx = adaptor.WithExecutor(context.Background(), db)
	Expect(t, session.InTx(ctx), BeFalse())
}
