package helper

import (
	"context"
	"database/sql"

	"github.com/xoctopus/sqlx/internal/sql/scanner"
	"github.com/xoctopus/sqlx/internal/structs"
	"github.com/xoctopus/sqlx/pkg/builder"
)

func ColumnsAndValuesForInsertion(m any) (builder.Cols, []any) {
	fields := structs.TableFields(m)
	cols := make([]builder.Col, 0, len(fields))
	vals := make([]any, 0, len(fields))
	for _, f := range fields {
		if !f.Field.ColumnDef.AutoInc && f.Value.IsValid() {
			cols = append(cols, builder.C(f.Field.ColumnName))
			vals = append(vals, f.Value.Interface())
		}
	}
	return builder.ColsOf(cols...), vals
}

func Scan(ctx context.Context, rows *sql.Rows, dst any) error {
	return scanner.Scan(ctx, rows, dst)
}

func AssignmentsOf() builder.Assignments {
	return nil
}
