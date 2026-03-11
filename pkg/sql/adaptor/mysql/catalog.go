package mysql

import (
	"context"
	"slices"
	"strings"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/sql/scanner"
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
	TableSchema       string  `db:"table_schema"`
	Table             string  `db:"table_name"`
	ColumnName        string  `db:"column_name"`
	RawDataType       string  `db:"data_type"`
	DataType          string  `db:"column_type"`
	VarcharLength     *uint64 `db:"character_maximum_length"` // char,varchar
	BinaryLength      *uint64 `db:"character_octet_length"`   // binary,varbinary
	DatetimePrecision *uint64 `db:"datetime_precision"`       // datetime/timestamp/time
	NumericWidth      *uint64 `db:"numeric_precision"`        // decimal,numeric,float,double
	NumericPrecision  *uint64 `db:"numeric_scale"`            // decimal,numeric
	DefaultValue      *string `db:"column_default"`
	IsNullable        string  `db:"is_nullable"`
	Comment           string  `db:"column_comment"`
	Extra             string  `db:"extra"`
	OrdinalPosition   string  `db:"ordinal_position"`
}

func (TSchemaTableColumn) TableName() string {
	return "information_schema.columns"
}

func (t *TSchemaTableColumn) ToCol() builder.Col {
	d := &def.ColumnDef{
		DefineFrom: def.DefineFromCatalog,
	}
	if len(t.Extra) > 0 {
		extra := strings.ToUpper(t.Extra)
		if extra == "AUTO_INCREMENT" {
			d.AutoInc = true
		} else {
			if after, ok := strings.CutPrefix(extra, "DEFAULT_GENERATED"); ok {
				following := after
				following = strings.TrimSpace(following)
				if after, ok := strings.CutPrefix(following, "ON UPDATE"); ok {
					onUpdate := after
					onUpdate = strings.TrimSpace(onUpdate)
					d.OnUpdate = new(onUpdate)
				}
			}
		}
	}
	if t.DefaultValue != nil {
		d.Default = t.DefaultValue
	}
	if strings.ToUpper(t.IsNullable) == "YES" {
		d.Null = true
	}
	d.DataType = strings.ToUpper(t.RawDataType)
	d.IsUnsigned = new(strings.HasSuffix(strings.ToUpper(t.DataType), "UNSIGNED"))
	if _, ok := TextBases[d.DataType]; ok {
		if t.VarcharLength != nil && *t.VarcharLength != TextBasesDefaultWidth[d.DataType] {
			d.Width = t.VarcharLength
		}
	}
	if _, ok := BinaryBases[d.DataType]; ok {
		if t.BinaryLength != nil && *t.BinaryLength != BinaryBasesDefaultWidth[d.DataType] {
			d.Width = t.BinaryLength
		}
	}
	if _, ok := DatetimeBases[d.DataType]; ok {
		if t.DatetimePrecision != nil && *t.DatetimePrecision != DatetimeBasesDefaultPrecision[d.DataType] {
			d.Precision = t.DatetimePrecision
		}
	}
	if _, ok := IntegerBases[d.DataType]; ok {
		if t.NumericWidth != nil {
			unsigned := *d.IsUnsigned
			if !unsigned && *t.NumericWidth != IntegerBasesDefaultWidth[d.DataType] ||
				unsigned && *t.NumericWidth != UnsignedIntegerBasesDefaultWidth[d.DataType] {
				d.Width = t.NumericWidth
			}
		}
	}
	if _, ok := FloatBases[d.DataType]; ok {
		if t.NumericPrecision != nil && t.NumericWidth != nil &&
			FloatBasesDefaultWidth[d.DataType] != *t.NumericWidth &&
			FloatBasesDefaultPrecision[d.DataType] != *t.NumericPrecision {
			d.Width = t.NumericWidth
			d.Precision = t.NumericPrecision
		}
	}
	d.Comment = t.Comment
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
			if name := strings.ToUpper(i.IndexName); name == "PRIMARY" {
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
