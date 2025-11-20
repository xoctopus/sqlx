package sqltypes_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
	_ "unsafe"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/sqltypes"
)

func TestDatetime_Hack(t *testing.T) {
	if os.Getenv("HACK_TEST") != "true" {
		t.Skip("should depend on postgres/mysql")
	}
	time.Sleep(10 * time.Second)

	var drivers = []struct {
		name string
		dsn  string
	}{
		{
			name: "postgres",
			dsn:  "postgres://postgres@localhost:15432/test?sslmode=disable",
		},
		{
			name: "postgres",
			dsn:  "postgres://postgres@localhost:15432/test?sslmode=disable&TimeZone=Asia/Shanghai",
		},
		{
			name: "mysql",
			dsn:  "root:@tcp(localhost:13306)/test",
		},
		{
			name: "mysql",
			dsn:  "root:@tcp(localhost:13306)/test?parseTime=true&loc=Asia%2FShanghai",
		},
	}
	sqltypes.AddTimestampInputLayouts(time.DateTime + ".000")

	query := "SELECT * FROM x LIMIT 1"

	for _, drv := range drivers {
		fmt.Printf("driver: %s %s\n", drv.name, drv.dsn)
		db := must.NoErrorV(sql.Open(drv.name, drv.dsn))

		rows, err := db.Query(query)
		if err != nil {
			t.Fatal(drv.name, err)
		}
		for rows.Next() {
			var (
				datetime  = &sqltypes.Datetime{}
				timestamp = &sqltypes.Timestamp{}
			)
			fmt.Println("scan:")
			err = rows.Scan(datetime, timestamp)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf(
				"datetime %s %d timestamp %s %d\n\n",
				datetime,
				datetime.Int(),
				timestamp,
				timestamp.Int(),
			)
		}
	}
}
