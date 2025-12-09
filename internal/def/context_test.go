package def_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	Expect(t, def.ModelTagKeyFrom(ctx), Equal("db"))

	ctx = def.WithModelTagKey(ctx, "gorm")
	Expect(t, def.ModelTagKeyFrom(ctx), Equal("gorm"))
}
