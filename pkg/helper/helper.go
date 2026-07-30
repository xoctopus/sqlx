package helper

import (
	"context"
	"database/sql"

	"github.com/xoctopus/x/reflectx"

	"github.com/xoctopus/sqlx/internal/structs"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/builder/modeled"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/sql/scanner"
)

// CVsForInsertion builds insert columns and flattened values from models.
// Auto-increment columns are skipped.
func CVsForInsertion[M builder.Model](ms ...M) (builder.Cols, []any) {
	if len(ms) == 0 {
		return nil, nil
	}

	m0 := (&modeled.Newer[M]{}).Model()
	fields := structs.TableFields(reflectx.IndirectNew(m0))

	cols := make([]builder.Col, 0, len(fields))
	for _, f := range fields {
		if !f.Field.ColumnDef.AutoInc && f.Value.IsValid() {
			cols = append(cols, builder.C(f.Field.ColumnName))
		}
	}

	vals := make([]any, 0, len(fields)*len(ms))
	for i := range ms {
		fs := structs.TableFields(reflectx.IndirectNew(&ms[i]))
		for _, f := range fs {
			if !f.Field.ColumnDef.AutoInc && f.Value.IsValid() {
				vals = append(vals, f.Value.Interface())
			}
		}
	}
	return builder.ColsOf(cols...), vals
}

// Scan scans rows into dst.
func Scan(ctx context.Context, rows *sql.Rows, dst any) error {
	return scanner.Scan(ctx, rows, dst)
}

// QueryAndScan executes f and scans the result into dst.
// If dst is nil, the query still runs and rows are discarded after close.
func QueryAndScan(ctx context.Context, e adaptor.Adaptor, f frag.Fragment, dst any) error {
	rows, err := e.Query(ctx, f)
	if err != nil {
		return err
	}

	defer func() {
		_ = rows.Close()
	}()

	if dst == nil {
		return nil
	}

	return scanner.Scan(ctx, rows, dst)
}
