package mysql

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type dialect struct{}

var _ adaptor.Dialect = (*dialect)(nil)

func (d dialect) CreateSchema(name string) frag.Fragment {
	return frag.Query("CREATE DATABASE IF NOT EXISTS ?;", frag.Lit(name))
}

func (d dialect) SwitchSchema(name string) frag.Fragment {
	return frag.Query("USE ?;", frag.Lit(name))
}

func (d dialect) CreateTableIfNotExists(t builder.Table) []frag.Fragment {
	exprs := []frag.Fragment{
		frag.Query(
			"CREATE TABLE IF NOT EXISTS @table (@def\n);", frag.NamedArgs{
				"table": t,
				"def": frag.Func(func(ctx context.Context) frag.Iter {
					return func(yield func(string, []any) bool) {
						i := 0
						for c := range t.Cols() {
							def := builder.GetColDef(c)
							if def.Deprecated != nil {
								continue
							}
							if i > 0 {
								yield(",", nil)
							}
							i++
							yield("\n\t", nil)

							for q, args := range c.Frag(ctx) {
								yield(q, args)
							}
							yield(" ", nil)
							for q, args := range d.DBType(def).Frag(ctx) {
								yield(q, args)
							}
						}
						for k := range t.Keys() {
							if k.IsPrimary() {
								f := frag.Query(
									",\n\tPRIMARY KEY (?)",
									builder.ColsIterOf(k.Cols()),
								)
								for q, args := range f.Frag(ctx) {
									yield(q, args)
								}
							}
						}
					}
				}),
			},
		),
	}
	for k := range t.Keys() {
		if !k.IsPrimary() {
			exprs = append(exprs, d.AddIndex(k))
		}
	}

	return exprs
}

func (d dialect) DropTable(t builder.Table) frag.Fragment {
	return frag.Query("DROP TABLE IF EXISTS @table;", frag.NamedArgs{"table": t})
}

func (d dialect) TruncateTable(t builder.Table) frag.Fragment {
	return frag.Query("TRUNCATE TABLE @table;", frag.NamedArgs{"table": t})
}

func (d dialect) AddColumn(c builder.Col) frag.Fragment {
	return frag.Query(
		"ALTER TABLE @table ADD COLUMN @col @datatype;",
		frag.NamedArgs{
			"table":    builder.GetColTable(c),
			"col":      c,
			"datatype": d.DBType(builder.GetColDef(c)),
		},
	)
}

func (d dialect) DropColumn(c builder.Col) frag.Fragment {
	return frag.Query(
		"ALTER TABLE @table DROP COLUMN @col;",
		frag.NamedArgs{
			"table": builder.GetColTable(c),
			"col":   c,
		},
	)
}

func (d dialect) RenameColumn(from builder.Col, to builder.Col) frag.Fragment {
	return frag.Query(
		"ALTER TABLE @table RENAME COLUMN @from TO @to;",
		frag.NamedArgs{
			"table": builder.GetColTable(from),
			"from":  from,
			"to":    to,
		},
	)
}

func (d dialect) ModifyColumn(curr, prev builder.Col) frag.Fragment {
	panic("todo")
}

func (d dialect) AddIndex(k builder.Key) frag.Fragment {
	if k.IsPrimary() {
		return frag.Query(
			"ALTER TABLE @table ADD PRIMARY KEY (@cols);", frag.NamedArgs{
				"table": k.(builder.WithTable).T(),
				"cols":  builder.ColsIterOf(k.Cols()),
			},
		)
	}

	def := k.(builder.KeyDef)
	return frag.Query(
		"CREATE @idx_type @idx_name ON @table (@cols)@idx_method;", frag.NamedArgs{
			"idx_name": frag.Lit(k.Name()),
			"idx_type": func() frag.Fragment {
				if k.IsUnique() {
					return frag.Lit("UNIQUE INDEX")
				}
				return frag.Lit("INDEX")
			}(),
			"table": k.(builder.WithTable).T(),
			"cols":  builder.KeyColumnsDefOf(k),
			"idx_method": func() frag.Fragment {
				if m := def.Method(); m != "" {
					return frag.Lit(" USING " + m)
				}
				return frag.Empty()
			}(),
		},
	)
}

func (d dialect) DropIndex(k builder.Key) frag.Fragment {
	tab := k.(builder.WithTable).T()

	// MUST remove auto_increment attribute first. but if it related other indexes?
	// cols := builder.ColsIterOf(k.Cols())
	// for c := range cols.Cols() {
	// 	def := builder.GetColDef(c)
	// 	if def.AutoInc {
	// 		// q = "ALTER TABLE @table MODIFY COLUMN @col @datatype;"
	// 		// args["col"] = c
	// 		// args["datatype"] = d.DBType(builder.GetColDef(c))
	// 		break
	// 	}
	// }

	if k.IsPrimary() {
		return frag.Query("ALTER TABLE ? DROP PRIMARY KEY;", tab)
	}
	return frag.Query("ALTER TABLE ? DROP INDEX ?;", tab, k)
}

func (d dialect) DBType(def builder.ColumnDef) frag.Fragment {
	modifiers := d.modifiers(def, d.datatype(def.Type, def))
	fragments := make([]frag.Fragment, 0, len(modifiers))

	for _, modifier := range modifiers {
		fragments = append(fragments, frag.Lit(modifier))
	}

	return frag.Compose(" ", fragments...)
}

func (d dialect) datatype(typ typx.Type, def builder.ColumnDef) string {
	// from catalog
	if def.DataType != "" {
		return def.DataType
	}

	must.BeTrueF(typ != nil, "column def missing type info")
	// from descriptor
	if rt, ok := typ.Unwrap().(reflect.Type); ok {
		v := reflect.New(rt).Interface()
		if desc, ok := v.(builder.WithDatatypeDesc); ok {
			return strings.ToUpper(desc.DBType("mysql"))
		}
	}

	datatype := ""
	switch kind := typ.Kind(); kind {
	case reflect.Pointer:
		datatype = d.datatype(typ.Elem(), def)
	case reflect.Bool:
		datatype = "BOOLEAN"
	case reflect.Int8:
		datatype = "TINYINT"
	case reflect.Uint8:
		datatype = "TINYINT UNSIGNED"
	case reflect.Int16:
		datatype = "SMALLINT"
	case reflect.Uint16:
		datatype = "SMALLINT UNSIGNED"
	case reflect.Int32, reflect.Int:
		datatype = "INT"
	case reflect.Uint32, reflect.Uint:
		datatype = "INT UNSIGNED"
	case reflect.Int64:
		datatype = "BIGINT"
	case reflect.Uint64:
		datatype = "BIGINT UNSIGNED"
	case reflect.Float32:
		datatype = "FLOAT"
	case reflect.Float64:
		datatype = "DOUBLE PRECISION"
	case reflect.String:
		if def.Width != 0 {
			datatype = "VARCHAR"
		} else {
			datatype = "TEXT"
		}
	default:
		if typ.PkgPath() == "time" && typ.Name() == "Time" {
			datatype = "DATETIME"
		} else {
			panic(fmt.Errorf("unsupported column type: %s", typ))
		}
	}
	return datatype
}

func (d dialect) modifiers(def builder.ColumnDef, datatype string) (modifiers []string) {
	// DataType => datatype(width,precision) eg: VARCHAR(width); DATETIME(precision); DECIMAL(width,precision);
	// here skip width determined types eg: BIGINT, TINYINT
	// ref: https://dev.mysql.com/doc/refman/8.0/en/numeric-type-attributes.html
	if (def.Width != 0 || def.Precision != 0) &&
		slices.Contains([]string{"VARCHAR", "DECIMAL", "NUMERIC", "DATETIME"}, datatype) {
		datatype += "("
		ss := make([]string, 0, 2)
		if def.Width != 0 {
			ss = append(ss, fmt.Sprintf("%d", def.Width))
		}
		if def.Precision != 0 {
			ss = append(ss, fmt.Sprintf("%d", def.Precision))
		}
		datatype += strings.Join(ss, ",")
		datatype += ")"
	}
	modifiers = append(modifiers, datatype)

	// Null ==> NOT NULL
	if !def.Null {
		modifiers = append(modifiers, "NOT NULL")
	}
	// Default ==> DEFAULT ...
	if v := def.Default; v != nil {
		modifiers = append(modifiers, "DEFAULT "+*v)
	}
	// OnUpdate ==> ON UPDATE ...
	if v := def.OnUpdate; v != nil {
		modifiers = append(modifiers, "ON UPDATE "+*v)
	}
	// AutoInc ==> AUTO_INCREMENT
	if def.AutoInc {
		modifiers = append(modifiers, "AUTO_INCREMENT")
	}
	// Comment ==> COMMENT '...'
	if v := def.Comment; v != "" {
		modifiers = append(modifiers, "COMMENT '"+v+"'")
	}
	return modifiers
}

func (d dialect) IsUnknownDatabaseError(err error) bool {
	return IsUnknownDatabaseError(err)
}

func (d dialect) IsConflictError(err error) bool {
	return IsConflictError(err)
}

func IsUnknownDatabaseError(err error) bool {
	var e *mysql.MySQLError
	return errors.As(err, &e) && e.Number == 1049
}

func IsConflictError(err error) bool {
	var e *mysql.MySQLError
	return errors.As(err, &e) && e.Number == 1062
}

func UnwrapError(err error) *mysql.MySQLError {
	var e *mysql.MySQLError
	if errors.As(err, &e) {
		return e
	}
	return nil
}
