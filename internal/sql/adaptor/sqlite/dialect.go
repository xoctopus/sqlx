package sqlite

/*
import (
	"errors"

	"modernc.org/sqlite"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type dialect struct {
}

func (d dialect) CreateSchema(string) frag.Fragment {
	panic("unsupported")
}

func (d dialect) SwitchSchema(string) frag.Fragment {
	panic("unsupported")
}

func (d dialect) CreateTableIfNotExists(t builder.Table) []frag.Fragment {
	panic("todo")
}

func (d dialect) DropTable(t builder.Table) frag.Fragment {
	return frag.Query("DROP TABLE IF EXISTS ?;", t)
}

func (d dialect) TruncateTable(t builder.Table) frag.Fragment {
	return frag.Query("TRUNCATE TABLE ?;", t)
}

func (d dialect) AddColumn(c builder.Col) frag.Fragment {
	return frag.Query("ALTER TABLE @table ADD COLUMN @column @datatype;", frag.NamedArgs{
		"table":    builder.GetColTable(c),
		"column":   c,
		"datatype": d.DBType(builder.GetColDef(c)),
	})
}

func (d dialect) DropColumn(c builder.Col) frag.Fragment {
	return frag.Query("ALTER TABLE @table DROP COLUMN @col;", frag.NamedArgs{
		"table": builder.GetColTable(c),
		"col":   c,
	})
}

func (d dialect) RenameColumn(from, to builder.Col) frag.Fragment {
	return frag.Query("ALTER TABLE @table RENAME COLUMN @from TO @to;", frag.NamedArgs{
		"table": builder.GetColTable(from),
		"from":  from,
		"to":    to,
	})
}

func (d dialect) ModifyColumn(curr, prev builder.Col) frag.Fragment {
	panic("todo")
}

func (d dialect) AddIndex(k builder.Key) frag.Fragment {
	if k.IsPrimary() {
		return nil
	}

	return frag.Query("CREATE @idx_type @idx_name ON @table (@options);", frag.NamedArgs{
		"table": builder.GetKeyTable(k),
		"idx_type": func() frag.Fragment {
			if k.IsUnique() {
				return frag.Lit("UNIQUE INDEX")
			}
			return frag.Lit("INDEX")
		}(),
		"idx_name": k.String(),
		"options":  builder.KeyColumnsDefOf(k),
	})
}

func (d dialect) DropIndex(k builder.Key) frag.Fragment {
	if k.IsPrimary() {
		return nil
	}
	return frag.Query("DROP INDEX IF EXISTS ?;", frag.Lit(k.String()))
}

func (d dialect) DBType(def builder.ColumnDef) frag.Fragment {
	panic("todo")
}

func (d dialect) IsUnknownDatabaseError(err error) bool {
	return false
}

func (d dialect) IsConflictError(err error) bool {
	var e *sqlite.Error
	return errors.As(err, &e) && e.Code() == 2067
}
*/
