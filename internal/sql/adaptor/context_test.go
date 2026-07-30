package adaptor

import (
	"context"
	"database/sql"
	"testing"

	. "github.com/xoctopus/x/testx"
)

type fakeExecutor struct {
	name string
}

func (fakeExecutor) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, nil
}

func (fakeExecutor) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, nil
}

func TestWithExecutor(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		Expect(t, ExecutorFrom(context.Background()), BeNil[Executor]())
	})

	t.Run("RoundTrip", func(t *testing.T) {
		e := fakeExecutor{name: "a"}
		ctx := WithExecutor(context.Background(), e)
		got := ExecutorFrom(ctx)
		Expect(t, got, NotBeNil[Executor]())
		Expect(t, got.(fakeExecutor).name, Equal("a"))
	})

	t.Run("Override", func(t *testing.T) {
		ctx := WithExecutor(context.Background(), fakeExecutor{name: "a"})
		ctx = WithExecutor(ctx, fakeExecutor{name: "b"})
		Expect(t, ExecutorFrom(ctx).(fakeExecutor).name, Equal("b"))
	})

	t.Run("NilExecutor", func(t *testing.T) {
		ctx := WithExecutor(context.Background(), nil)
		Expect(t, ExecutorFrom(ctx), BeNil[Executor]())
	})
}
