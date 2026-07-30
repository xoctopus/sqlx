package mysql

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/def"
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
								if !yield(",", nil) {
									return
								}
							}
							i++
							if !yield("\n\t", nil) {
								return
							}

							for q, args := range c.Frag(ctx) {
								if !yield(q, args) {
									return
								}
							}
							yield(" ", nil)
							for q, args := range d.DBType(def).Frag(ctx) {
								if !yield(q, args) {
									return
								}
							}
						}
						for k := range t.Keys() {
							if k.IsPrimary() {
								f := frag.Query(
									",\n\tPRIMARY KEY (?)",
									builder.ColsIterOf(k.Cols()),
								)
								for q, args := range f.Frag(ctx) {
									if !yield(q, args) {
										return
									}
								}
							}
						}
					}
				}),
			},
		),
	}

	var keys []builder.Key
	for k := range t.Keys() {
		if !k.IsPrimary() {
			keys = append(keys, k)
		}
	}
	slices.SortFunc(keys, func(a, b builder.Key) int {
		if a.IsUnique() && !b.IsUnique() {
			return -1
		}
		if !a.IsUnique() && b.IsUnique() {
			return 1
		}
		return cmp.Compare(a.Name(), b.Name())
	})
	for _, k := range keys {
		exprs = append(exprs, d.AddIndex(k))
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

func (d dialect) ModifyColumn(next, curr builder.Col) frag.Fragment {
	nextDef := builder.GetColDef(next)
	currDef := builder.GetColDef(curr)

	if nextDef.AutoInc {
		// maybe multi steps
		return nil
	}

	typNext, _ := frag.Collect(context.Background(), d.DBType(nextDef))
	typCurr, _ := frag.Collect(context.Background(), d.DBType(currDef))

	if strings.EqualFold(typNext, typCurr) {
		return nil
	}

	return frag.Query(
		"ALTER TABLE @table MODIFY COLUMN @col @next; -- from @prev",
		frag.NamedArgs{
			"table": builder.GetColTable(next),
			"col":   next,
			"next":  frag.Lit(typNext),
			"prev":  frag.Lit(typCurr),
		},
	)
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
	d.datatype(def.Type, &def)
	modifiers := d.modifiers(def)
	fragments := make([]frag.Fragment, 0, len(modifiers))

	for _, modifier := range modifiers {
		fragments = append(fragments, frag.Lit(modifier))
	}

	return frag.Compose(" ", fragments...)
}

func (d dialect) ColDefine(dd *builder.ColumnDef) {
	must.BeTrueF(
		dd != nil && dd.Type != nil,
		"invalid column define",
	)
	d.datatype(dd.Type, dd)
}

func (d dialect) datatype(typ typx.Type, dd *builder.ColumnDef) {
	switch dd.DefineFrom {
	case def.DefineFromCatalog, def.DefineFromUser:
		return
	}

	must.BeTrueF(typ != nil, "column def missing type info")
	if rt, ok := typ.Unwrap().(reflect.Type); ok {
		rv := reflect.New(rt)
		ptr := rv.Interface()
		if desc, ok := ptr.(builder.WithDatatypeDesc); ok {
			dd.DataType = strings.ToUpper(desc.DBType("mysql"))
			return
		} else {
			val := rv.Elem().Interface()
			if desc, ok = val.(builder.WithDatatypeDesc); ok {
				dd.DataType = strings.ToUpper(desc.DBType("mysql"))
				return
			}
		}
	}

	if dd.DataType == "" {
		switch kind := typ.Kind(); kind {
		case reflect.Pointer:
			d.datatype(typ.Elem(), dd)
		case reflect.Bool:
			dd.DataType = "TINYINT"
		case reflect.Int8, reflect.Uint8:
			dd.DataType = "TINYINT"
		case reflect.Int16, reflect.Uint16:
			dd.DataType = "SMALLINT"
		case reflect.Int32, reflect.Int, reflect.Uint32, reflect.Uint:
			dd.DataType = "INT"
		case reflect.Int64, reflect.Uint64:
			dd.DataType = "BIGINT"
		case reflect.Float32:
			dd.DataType = "FLOAT"
		case reflect.Float64:
			dd.DataType = "DOUBLE PRECISION"
		case reflect.String:
			if dd.Width != nil && *dd.Width != 0 {
				dd.DataType = "VARCHAR"
			} else {
				dd.DataType = "TEXT"
			}
		default:
			if typ.PkgPath() == "time" && typ.Name() == "Time" {
				dd.DataType = "DATETIME"
			} else {
				panic(fmt.Errorf("unsupported column type: %s", typ))
			}
		}
	}

	if dd.IsUnsigned == nil {
		switch typ.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dd.IsUnsigned = new(true)
		default:
		}
	}
}

func (d dialect) modifiers(dd builder.ColumnDef) (modifiers []string) {
	// DataType => datatype(width,precision) eg: VARCHAR(width); DATETIME(precision); DECIMAL(width,precision);
	// here skip width determined types eg: BIGINT, TINYINT, FLOAT, DOUBLE
	// ref: https://dev.mysql.com/doc/refman/8.0/en/numeric-type-attributes.html
	datatype := dd.DataType
	s := ""
	if dd.Width != nil || dd.Precision != nil {
		if dd.Width != nil && *dd.Width > 0 {
			s += strconv.FormatUint(*dd.Width, 10)
		}
		if dd.Precision != nil && *dd.Precision > 0 {
			if dd.Width != nil {
				s += ","
			}
			s += strconv.FormatUint(*dd.Precision, 10)
		}
	}
	if len(s) > 0 {
		kind := reflect.Invalid
		if dd.Type != nil {
			kind = dd.Type.Kind()
		}
		switch kind {
		case reflect.Float32, reflect.Float64:
		default:
			datatype += "(" + s + ")"
		}
	}

	unsigned := dd.IsUnsigned != nil && *dd.IsUnsigned
	if v, ok := DefaultWithWidth(dd.DataType, unsigned); ok && v == datatype {
		datatype = dd.DataType
	}

	modifiers = append(modifiers, datatype)

	if dd.IsUnsigned != nil && *dd.IsUnsigned {
		modifiers = append(modifiers, "UNSIGNED")
	}

	// Null ==> NOT NULL
	if !dd.Null {
		modifiers = append(modifiers, "NOT NULL")
	}
	// Default ==> DEFAULT ...
	if v := dd.Default; v != nil {
		p := *v
		needquoted := false
		if _, ok := MustQuoteBases[dd.DataType]; ok {
			needquoted = true
			up := strings.ToUpper(p)
			for _, k := range DefaultKeywords { // Keywords like NOW CURRENT_TIMESTAMP
				if strings.HasPrefix(up, k) {
					needquoted = false
					break
				}
			}
			if needquoted {
				if dd.DefineFrom != def.DefineFromCatalog {
					needquoted = !strings.HasPrefix(p, "'")
				}
			}

			if needquoted {
				needquoted = !(strings.HasPrefix(p, "(") && strings.HasSuffix(p, ")"))
			}
		}
		if needquoted {
			p = "'" + strings.ReplaceAll(p, "'", "''") + "'"
		}

		modifiers = append(modifiers, "DEFAULT")
		modifiers = append(modifiers, p)
	}
	// OnUpdate ==> ON UPDATE ...
	if v := dd.OnUpdate; v != nil {
		modifiers = append(modifiers, "ON UPDATE "+*v)
	}
	// AutoInc ==> AUTO_INCREMENT
	if dd.AutoInc {
		modifiers = append(modifiers, "AUTO_INCREMENT")
	}
	// // Comment ==> COMMENT '...'
	// if v := dd.Comment; v != "" {
	// 	modifiers = append(modifiers, "COMMENT '"+v+"'")
	// }
	return modifiers
}

func (d dialect) IsUnknownDatabaseError(err error) bool {
	return IsUnknownDatabaseError(err)
}

func (d dialect) IsConflictError(err error) bool {
	return IsConflictError(err)
}

func IsUnknownDatabaseError(err error) bool {
	target, ok := errors.AsType[*mysql.MySQLError](err)
	return ok && target.Number == 1049
}

func IsConflictError(err error) bool {
	target, ok := errors.AsType[*mysql.MySQLError](err)
	return ok && target.Number == 1062
}

func UnwrapError(err error) *mysql.MySQLError {
	if target, ok := errors.AsType[*mysql.MySQLError](err); ok {
		return target
	}
	return nil
}

func ErrorLevel(err error) int {
	if IsConflictError(err) {
		return 0
	}
	return 1
}
