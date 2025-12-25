package migrator

import (
	"context"
	"database/sql/driver"
	"fmt"
	"slices"

	"github.com/xoctopus/sqlx/internal/diff"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

var CtxMode = diff.CtxMode

// Migrate migrates tables
// TODOs:
// 1. check use create table mode only
// 2. if adaptor supports DDL transaction
func Migrate(ctx context.Context, a adaptor.Adaptor, next builder.Catalog) (string, error) {
	curr, err := a.Catalog(ctx)
	if err != nil {
		return "", err
	}

	fragments := make([]frag.Fragment, 0)

	for _, name := range slices.Sorted(builder.TableNames(next)) {
		d := diff.Diff(ctx, a.Dialect(), curr.T(name), next.T(name))
		if frag.IsNil(d) {
			continue
		}
		fragments = append(fragments, d)
	}

	if len(fragments) == 0 {
		return "", nil
	}

	q, args := frag.Collect(ctx, frag.Compose("\n", fragments...))
	named := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		named[i].Value = arg
	}

	q, _ = loggingdriver.DefaultInterpolate(q, named)

	if mode, ok := CtxMode.From(ctx); ok && mode.Is(diff.MODE_DRY_RUN) {
		return q, nil
	}

	return q, a.Tx(ctx, func(ctx context.Context) error {
		for _, m := range fragments {
			if _, err := a.Exec(ctx, m); err != nil {
				return fmt.Errorf("migrate failed: %w", err)
			}
		}
		return nil
	})
}
