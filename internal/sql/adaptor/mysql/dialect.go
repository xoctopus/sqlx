package mysql

import (
	"context"
	"fmt"
	"reflect"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type dialect struct{}

var _ adaptor.Dialect = (*dialect)(nil)

func (d *dialect) CreateSchema(name string) frag.Fragment {
	return nil
}

func (d *dialect) SwitchSchema(name string) frag.Fragment {
	return frag.Query("USE ?;", name)
}

/*
CREATE TABLE IF NOT EXISTS users (
    f_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    f_user_id BIGINT UNSIGNED NOT NULL,
    f_name VARCHAR(64) not null,
    f_email VARCHAR(255) NOT NULL,
    f_created_at BIGINT UNSIGNED NOT NULL,
    f_updated_at BIGINT UNSIGNED NOT NULL,
    f_deleted_at BIGINT UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

    UNIQUE KEY ui_email (f_email,f_deleted_at),
    INDEX i_name (f_name,f_deleted_at),
    INDEX i_created_at (f_created_at),
    INDEX i_updated_at (f_updated_at),
    UNIQUE KEY ui_user_id (f_user_id,f_deleted_at),
    PRIMARY KEY (f_id)
*/

func (d *dialect) CreateTableIsNotExists(t builder.Table) []frag.Fragment {
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

func (d *dialect) DropTable(t builder.Table) frag.Fragment {
	return frag.Query("DROP TABLE IF EXISTS @table;", frag.NamedArgs{"table": t})
}

func (d *dialect) TruncateTable(t builder.Table) frag.Fragment {
	return frag.Query("TRUNCATE TABLE @table;", frag.NamedArgs{"table": t})
}

func (d *dialect) AddColumn(c builder.Col) frag.Fragment {
	return frag.Query(
		"ALTER TABLE @table ADD COLUMN @col @datatype;",
		frag.NamedArgs{
			"table":    builder.GetColTable(c),
			"col":      c,
			"datatype": d.DBType(builder.GetColDef(c)),
		},
	)
}

func (d *dialect) DropColumn(c builder.Col) frag.Fragment {
	return frag.Query(
		"ALTER TABLE @table DROP COLUMN @col;",
		frag.NamedArgs{
			"table": builder.GetColTable(c),
			"col":   c,
		},
	)
}

func (d *dialect) RenameColumn(from builder.Col, to builder.Col) frag.Fragment {
	return frag.Query(
		"ALTER TABLE @table RENAME COLUMN @from TO @to;",
		frag.NamedArgs{
			"table": builder.GetColTable(from),
			"from":  from,
			"to":    to,
		},
	)
}

func (d *dialect) ModifyColumn(builder.Col, builder.Col) frag.Fragment {
	return nil
}

func (d *dialect) AddIndex(k builder.Key) frag.Fragment {
	if k.IsPrimary() {
		return frag.Query(
			"ALTER TABLE @table ADD PRIMARY KEY (@cols)", frag.NamedArgs{
				"table": k.(builder.WithTable).T(),
				"cols":  builder.ColsIterOf(k.Cols()),
			},
		)
	}

	def := k.(builder.KeyDef)
	return frag.Query(
		"CREATE @idx_type @idx_name ON @table (@cols)@idx_method;", frag.NamedArgs{
			"idx_name": k.Name(),
			"idx_type": func() frag.Fragment {
				if k.IsUnique() {
					return frag.Lit("UNIQUE INDEX")
				}
				return frag.Lit("INDEX")
			},
			"@table": k.(builder.WithTable).T,
			"@cols":  builder.KeyColumnsDefOf(k),
			"@idx_method": func() frag.Fragment {
				if m := def.Method(); m != "" {
					return frag.Lit(" USING " + m)
				}
				return frag.Empty()
			},
		},
	)
}

func (d *dialect) DropIndex(k builder.Key) frag.Fragment {
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

func (d *dialect) DBType(def builder.ColumnDef) frag.Fragment {
	return nil
}

func (d *dialect) datatype(typ typx.Type, def builder.ColumnDef) string {
	// from catalog
	if def.DataType != "" {
		return def.DataType
	}

	must.BeTrueF(typ != nil, "column def missing type info")
	// from descriptor
	if rt, ok := typ.Unwrap().(reflect.Type); ok {
		v := reflect.New(rt).Interface()
		if desc, ok := v.(builder.WithDatatypeDesc); ok {
			return desc.DBType("mysql")
		}
	}

	switch kind := typ.Kind(); kind {
	case reflect.Pointer:
		return d.datatype(typ.Elem(), def)
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int8:
		return "TINYINT"
	case reflect.Uint8:
		return "TINYINT UNSIGNED"
	case reflect.Int16:
		return "SMALLINT"
	case reflect.Uint16:
		return "SMALLINT UNSIGNED"
	case reflect.Int32:
		return "INT"
	case reflect.Uint32:
		return "INT UNSIGNED"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint64:
		return "BIGINT UNSIGNED"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.String:
		return "TEXT"
	default:
		if typ.PkgPath() == "time" && typ.Name() == "Time" {
			return "DATETIME"
		}
		panic(fmt.Errorf("unsupported column type: %s", typ))
	}
}

func (d *dialect) modifiers(def builder.ColumnDef, t string) (modifiers []string) {
	if !def.Null {
		modifiers = append(modifiers, "NOT NULL")
	}
	if v := def.Default; v != nil {
		modifiers = append(modifiers, "DEFAULT '"+*v+"'")
	}
	return modifiers
}
