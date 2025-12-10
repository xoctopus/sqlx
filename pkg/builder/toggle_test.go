package builder_test

import (
	"context"
	"testing"

	"github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
)

func TestToggles(t *testing.T) {
	ts := builder.TogglesFromContext(nil)

	testx.Expect(t, len(ts), testx.Equal(0))

	ctx := builder.WithToggles(context.Background(), 1, 2)
	testx.Expect(t, builder.HasToggle(ctx, 1), testx.BeTrue())
	testx.Expect(t, builder.HasToggle(ctx, 2), testx.BeTrue())
	testx.Expect(t, builder.HasToggle(ctx, 3), testx.BeFalse())
}
