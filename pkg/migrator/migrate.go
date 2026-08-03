package migrator

import (
	"context"
	"database/sql/driver"
	"fmt"
	"slices"

	"github.com/xoctopus/x/contextx"

	"github.com/xoctopus/sqlx/internal/diff"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/migrator/internal"
)

var (
	// CtxMode carries migration/diff mode flags.
	CtxMode = diff.CtxMode
	// CtxMeta carries migration sql meta switcher
	CtxMeta = contextx.NewV[bool](true)
	// CtxOutput is reserved for writing migration output to a directory.
	CtxOutput = contextx.NewT[string]()
)

// Mode is a migration/diff mode flag set.
type Mode = diff.Mode

const (
	// DIFF_MODE_CREATE_TABLE skips alter when the current table already exists.
	DIFF_MODE_CREATE_TABLE = diff.MODE_CREATE_TABLE
	// DIFF_MODE_DRY_RUN returns SQL without executing it.
	DIFF_MODE_DRY_RUN = diff.MODE_DRY_RUN
)

// Migrate diffs curr catalog against next, appends SQL meta document statements,
// and executes them in a transaction unless DIFF_MODE_DRY_RUN is set.
// It returns the interpolated SQL script.
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

	if enabled, _ := CtxMeta.From(ctx); enabled {
		fragments = append(
			fragments,
			slices.Concat(
				internal.GenerateTableDocuments(ctx, a, next),
				internal.GenerateTableColumnDocuments(ctx, a, next),
				internal.GenerateTableEnumerationDocument(ctx, a, next),
			)...,
		)
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
