package sqlx

import "strings"

type Attr string

const (
	// AttrTableName defines model's table name to extend TableName method
	AttrTableName Attr = "TableName"
	// AttrRegister defines model's catalog to register table
	AttrRegister Attr = "Register"
)

var attrs = []Attr{
	AttrTableName,
	AttrRegister,
}

func HasAttr(x string) (a Attr) {
	for _, v := range attrs {
		if strings.EqualFold(string(v), x) {
			a = v
			break
		}
	}
	return a
}
