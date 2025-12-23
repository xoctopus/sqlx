package mysql

import (
	"context"
	"slices"
	"strings"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/scanner"
	"github.com/xoctopus/sqlx/pkg/builder"
)

type TSchemaTableIndex struct {
	TableSchema string `db:"table_schema"`
	IndexType   string `db:"index_type"`
	IndexName   string `db:"index_name"`
	NonUnique   int    `db:"non_unique"`
	Table       string `db:"table_name"`
	ColumnName  string `db:"column_name"`
	SeqInIndex  int    `db:"seq_in_index"`
}

func (TSchemaTableIndex) TableName() string {
	return "information_schema.statistics"
}

type TSchemaTableColumn struct {
	TableSchema       string `db:"table_schema"`
	Table             string `db:"table_name"`
	ColumnName        string `db:"column_name"`
	RawDataType       string `db:"data_type"`
	DataType          string `db:"column_type"`
	VarcharLength     uint64 `db:"character_maximum_length"` // char,varchar
	BinaryLength      uint64 `db:"character_octet_length"`   // binary,varbinary
	NumericWidth      uint64 `db:"numeric_precision"`        // decimal,numeric,float,double
	NumericPrecision  uint64 `db:"numeric_scale"`            // decimal,numeric
	DatetimePrecision uint64 `db:"datetime_precision"`       // datetime/timestamp/time
	DefaultValue      string `db:"column_default"`
	IsNullable        string `db:"is_nullable"`
	Extra             string `db:"extra"`
	OrdinalPosition   string `db:"ordinal_position"`
}

func (TSchemaTableColumn) TableName() string {
	return "information_schema.columns"
}

func (t *TSchemaTableColumn) ToCol() builder.Col {
	d := &def.ColumnDef{}
	if strings.ToLower(t.Extra) == "auto_increment" {
		d.AutoInc = true
	}
	if t.DefaultValue != "" {
		d.Default = &t.DefaultValue
	}
	if t.IsNullable == "YES" {
		d.Null = true
	}
	d.DataType = t.DataType
	datatype := strings.ToLower(t.RawDataType)
	switch datatype {
	case "char", "varchar":
		d.Width = t.VarcharLength
	case "binary", "varbinary":
		d.Width = t.BinaryLength
	case "datetime", "timestamp", "time":
		d.Precision = t.DatetimePrecision
	case "decimal", "numeric":
		// skip float/double precision width and precision define
		// https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html
		d.Width = t.NumericPrecision
		d.Precision = t.NumericPrecision
	}
	return builder.C(t.ColumnName, builder.WithColDef(d))
}

func ScanCatalog(ctx context.Context, a adaptor.Adaptor, database string) (builder.Catalog, error) {
	catalog := builder.NewCatalog()

	tC := builder.TFrom(&TSchemaTableColumn{})
	expr := builder.Select(builder.ColsIterOf(tC.Cols())).
		From(
			tC,
			builder.Where(
				builder.CC[string](tC.C("table_schema")).AsCond(builder.Eq(database)),
			),
			builder.OrderBy(
				builder.Order(tC.C("table_name")),
				builder.Order(tC.C("ordinal_position")),
			),
		)
	rows, err := a.Query(ctx, expr)
	if err != nil {
		return nil, err
	}
	columns := make([]*TSchemaTableColumn, 0)
	if err = scanner.Scan(ctx, rows, &columns); err != nil {
		return nil, err
	}
	for _, s := range columns {
		var t builder.Table
		if t = catalog.T(s.Table); t == nil {
			t = builder.T(s.Table)
			t = t.(builder.WithSchema).WithSchema(a.Schema())
			t = t.(builder.WithDatabase).WithDatabase(a.Endpoint())
			catalog.Add(t)
		}
		t.(builder.ColsManager).AddCol(s.ToCol())
	}

	tI := builder.TFrom(&TSchemaTableIndex{})
	expr = builder.Select(builder.ColsIterOf(tI.Cols())).
		From(
			tI,
			builder.Where(
				builder.CC[string](tC.C("table_schema")).AsCond(builder.Eq(database)),
			),
			builder.OrderBy(
				builder.Order(tI.C("table_name")),
				builder.Order(tI.C("index_name")),
				builder.Order(tI.C("seq_in_index")),
			),
		)
	rows, err = a.Query(ctx, expr)
	if err != nil {
		return nil, err
	}
	indexes := make([]*TSchemaTableIndex, 0)
	if err = scanner.Scan(ctx, rows, &indexes); err != nil {
		return nil, err
	}
	grouped := make(map[string]map[string][]*TSchemaTableIndex)
	for _, i := range indexes {
		if grouped[i.Table] == nil {
			grouped[i.Table] = make(map[string][]*TSchemaTableIndex)
		}
		grouped[i.Table][i.IndexName] = append(grouped[i.Table][i.IndexName], i)
	}
	for table := range grouped {
		for index := range grouped[table] {
			t := catalog.T(table)
			must.BeTrueF(t != nil, "table %s not scanned from information_schema.columns", table)

			cols := make([]builder.Col, 0)
			list := grouped[table][index]
			slices.SortFunc(list, func(a, b *TSchemaTableIndex) int {
				return a.SeqInIndex - b.SeqInIndex
			})
			for _, i := range list {
				cols = append(cols, t.C(i.ColumnName))
			}

			i := list[0]
			options := make([]builder.KeyOption, 0) // skip empty and default
			if i.IndexType != "" {
				options = append(options, builder.WithKeyMethod(i.IndexType))
			}

			k := builder.Key(nil)
			if name := strings.ToLower(i.IndexName); name == "primary" {
				k = builder.PK(builder.ColsOf(cols...), options...)
			} else {
				if i.NonUnique != 0 {
					k = builder.K(name, builder.ColsOf(cols...), options...)
				} else {
					k = builder.UK(name, builder.ColsOf(cols...), options...)
				}
			}

			t.(builder.KeysManager).AddKey(k)
		}
	}
	return catalog, nil
}
