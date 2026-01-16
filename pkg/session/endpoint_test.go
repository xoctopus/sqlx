package session_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/pkg/session"
)

func TestEndpoint(t *testing.T) {
	ctx := hack.WithSession(
		hack.Context(t),
		t,
		"mysql://root@localhost:13306/test",
	)
	sess := session.From(ctx, "sqlx.hack")
	Expect(t, sess, NotBeNil[session.Session]())
}
