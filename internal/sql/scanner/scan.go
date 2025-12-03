package scanner

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/reflectx"

	"github.com/xoctopus/sqlx/internal/sql/scanner/nullable"
	"github.com/xoctopus/sqlx/internal/structs"
	sqlerrs "github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func Scan(ctx context.Context, rows *sql.Rows, v any) (err error) {
	if rows == nil {
		return nil
	}

	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var si ScanIter

	si, err = ScanIterFor(v)
	if err == nil {
		for rows.Next() {
			item := si.New()
			if err = scan(ctx, rows, item); err != nil {
				return err
			}
			_ = si.Next(item)
		}

		if x, ok := si.(interface{ MustHasRecord() bool }); ok {
			if !x.MustHasRecord() {
				return codex.Errorf(sqlerrs.NOTFOUND, "record is not found")
			}
		}

		if err = rows.Err(); err != nil {
			return err
		}
	}

	return
}

func scan(ctx context.Context, rows *sql.Rows, v any) error {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("scan target must be a ptr value, but got %T", v)
	}
	if s, ok := v.(sql.Scanner); ok {
		return rows.Scan(s)
	}
	t = reflectx.Deref(t)

	switch t.Kind() {
	case reflect.Struct:
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		n := len(cols)
		if n < 1 {
			return nil
		}

		dst := make([]any, n)
		placeholder := &nullable.EmptyScanner{}

		if x, ok := v.(WithColumnReceivers); ok {
			receivers := x.ColumnReceivers()
			for i, name := range cols {
				if c, ok := receivers[strings.ToLower(name)]; ok {
					dst[i] = nullable.NewNullIgnoreScanner(c)
				} else {
					dst[i] = &placeholder
				}
			}
			return rows.Scan(dst...)
		}

		columns := map[string]int{}

		for i, name := range cols {
			columns[strings.ToLower(name)] = i
			dst[i] = placeholder
		}

		for _, f := range structs.TableFields(ctx, v) {
			if f.TableName != "" {
				if i, ok := columns[frag.Alias(f.TableName, f.Field.Name)]; ok && i > -1 {
					dst[i] = nullable.NewNullIgnoreScanner(f.Value.Addr().Interface())
				}
			}
			if i, ok := columns[f.Field.Name]; ok && i > -1 {
				dst[i] = nullable.NewNullIgnoreScanner(f.Value.Addr().Interface())
			}
		}
		return rows.Scan(dst...)
	default:
		return rows.Scan(nullable.NewNullIgnoreScanner(v))
	}
}
