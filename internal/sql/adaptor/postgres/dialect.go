package postgres

/*
import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type dialect struct{}

func (d dialect) CreateSchema(name string) frag.Fragment {
	return frag.Query("CREATE SCHEMA IF NOT EXISTS ?;", frag.Lit(name))
}

func (d dialect) SwitchSchema(name string) frag.Fragment {
	return frag.Query("SET search_path TO ?;", frag.Lit(name))
}

func (d dialect) CreateTableIfNotExists(t builder.Table) []frag.Fragment {
	panic("todo")
}

func (d dialect) DropTable(t builder.Table) frag.Fragment {
	return frag.Query("DROP TABLE IF EXISTS @table;", frag.NamedArgs{
		"table": t,
	})
}

func (d dialect) TruncateTable(t builder.Table) frag.Fragment {
	return frag.Query("TRUNCATE TABLE @table;", frag.NamedArgs{
		"table": t,
	})
}

func (d dialect) AddColumn(c builder.Col) frag.Fragment {
	return frag.Query("ALTER TABLE @table ADD COLUMN @column @datatype;", frag.NamedArgs{
		"table":    builder.GetColTable(c),
		"column":   c,
		"datatype": d.DBType(builder.GetColDef(c)),
	})
}

func (d dialect) DropColumn(c builder.Col) frag.Fragment {
	return frag.Query("ALTER TABLE @table DROP COLUMN @column;", frag.NamedArgs{
		"table":  builder.GetColTable(c),
		"column": c,
	})
}

func (d dialect) RenameColumn(from builder.Col, to builder.Col) frag.Fragment {
	return frag.Query("ALTER TABLE @table RENAME COLUMN @from TO @to;", frag.NamedArgs{
		"table": builder.GetColTable(from),
		"from":  builder.GetColTable(from),
		"to":    builder.GetColTable(to),
	})
}

func (d dialect) ModifyColumn(curr, prev builder.Col) frag.Fragment {
	panic("todo")
}

func (d dialect) AddIndex(k builder.Key) frag.Fragment {
	panic("todo")
}

func (d dialect) DropIndex(k builder.Key) frag.Fragment {
	name := k.String()
	if k.IsPrimary() {
		name = "pkey"
	}

	if k.IsPrimary() {
		return frag.Query("ALTER TABLE @table DROP CONSTRAINT @index;", frag.NamedArgs{
			"table": builder.GetKeyTable(k),
			"index": name,
		})
	}
	return frag.Query("DROP INDEX IF EXISTS @index?", frag.NamedArgs{"index": name})
}

func (d dialect) DBType(def builder.ColumnDef) frag.Fragment {
	panic("todo")
}

func (d dialect) IsUnknownDatabaseError(err error) bool {
	var e *pgconn.PgError
	return errors.As(err, &e) && e.Code == "3D000"
}

func (d dialect) IsConflictError(err error) bool {
	var e *pgconn.PgError
	return errors.As(err, &e) && e.Code == "23505"
}
*/
