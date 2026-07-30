package testdata

import "github.com/xoctopus/sqlx/pkg/types"

// TestTable 测试用
// +genx:model
// @model TableName=t_test_table
type TestTable struct {
	Int8  int8  `db:"f_int8"`
	Int16 int16 `db:"f_int16"`
	Int32 int32 `db:"f_int32"`
	Int64 int64 `db:"f_int64"`

	Uint8  uint8  `db:"f_uint8"`
	Uint16 uint16 `db:"f_uint16"`
	Uint32 uint32 `db:"f_uint32"`
	Uint64 uint64 `db:"f_uint64"`

	Int8Default int8 `db:"f_int8_default,default=10"`

	String_          string `db:"f_string"`
	StringWidthFixed string `db:"f_string_w,width=128"`

	Float32        float32 `db:"f_float32,precision=4"`
	Float32Default float32 `db:"f_float32_default,precision=4,default=1.2"`

	Decimal        types.Decimal `db:"f_decimal"`
	DecimalDefault types.Decimal `db:"f_decimal_default,width=32,precision=5,default=1.20000"`
}
