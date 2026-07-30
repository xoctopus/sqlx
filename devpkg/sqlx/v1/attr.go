package sqlx

type Attr string

const (
	// AttrTableName defines model's table name to extend TableName method
	AttrTableName Attr = "TableName"
	// AttrRegister defines model's catalog to register table
	AttrRegister Attr = "Register"
)
