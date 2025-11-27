package types_test

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

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestDatetime_Hack(t *testing.T) {
	if os.Getenv("HACK_TEST") != "true" {
		t.Skip("should depend on postgres/mysql")
	}
	time.Sleep(10 * time.Second)

	types.SetTimestampPrecision(types.TIMESTAMP_PRECISION__MILLI)
	types.SetTimestampTimezone(types.CST)

	fmt.Println("parse precision: ", types.GetTimestampPrecision())
	fmt.Println("parse timezone:  ", types.GetTimestampTimezone())
	fmt.Println("==========================")

	var drivers = []struct {
		driver string
		label  string
		dsn    string
	}{
		// {
		// 	name: "postgres",
		// 	dsn:  "postgres://postgres@localhost:15432/test?sslmode=disable",
		// },
		// {
		// 	name: "postgres",
		// 	dsn:  "postgres://postgres@localhost:15432/test?sslmode=disable&TimeZone=Asia/Shanghai",
		// },
		{
			driver: "mysql",
			label:  "mysql",
			dsn:    "root:@tcp(localhost:13306)/test",
		},
		{
			driver: "mysql",
			label:  "mysql-with-cst",
			dsn:    "root:@tcp(localhost:13306)/test?loc=Asia%2FShanghai",
		},
		{
			driver: "mysql",
			label:  "mysql-with-jst",
			dsn:    "root:@tcp(localhost:13306)/test?loc=Asia%2FTokyo",
		},
		{
			driver: "mysql",
			label:  "mysql-with-cst-and-time-parsing",
			dsn:    "root:@tcp(localhost:13306)/test?parseTime=true&loc=Asia%2FShanghai",
		},
		{
			driver: "mysql",
			label:  "mysql-with-jst-and-time-parsing",
			dsn:    "root:@tcp(localhost:13306)/test?parseTime=true&loc=Asia%2FTokyo",
		},
	}
	types.AddTimestampInputLayouts(time.DateTime + ".000")

	query := "SELECT * FROM x LIMIT 1"

	for _, drv := range drivers {
		fmt.Printf("driver: %s\n", drv.driver)
		fmt.Printf("label:  %s\n", drv.label)
		fmt.Printf("dsn:    %s\n", drv.dsn)

		db := must.NoErrorV(sql.Open(drv.driver, drv.dsn))
		rows, err := db.Query(`SELECT
		@@system_time_zone  AS system_zone,  -- os timezone
		@@global.time_zone  AS global_zone,  -- mysql global timezone
		@@session.time_zone AS session_zone; -- current connection`)
		if err != nil {
			t.Fatal(drv.label, err)
		}
		for rows.Next() {
			var (
				system  string
				global  string
				session string
			)
			err = rows.Scan(&system, &global, &session)
			if err != nil {
				t.Fatal(drv.label, err)
			}
			fmt.Printf(
				"driver timezone info %s: [system: %s] [global: %s] [session: %s]\n",
				drv.label, system, global, session,
			)
		}

		rows, err = db.Query(query)
		if err != nil {
			t.Fatal(drv.label, err)
		}
		for rows.Next() {
			var (
				datetime  = &types.Datetime{}
				timestamp = &types.Timestamp{}
			)
			fmt.Print("scan ")
			err = rows.Scan(datetime, timestamp)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("\nresult:")
			fmt.Printf("datetime:  %s %d\n", datetime.Format(time.RFC3339Nano), datetime.UnixNano())
			fmt.Printf("timestamp: %s %d\n", timestamp.Format(time.RFC3339Nano), timestamp.UnixNano())
		}
		db.Close()
		fmt.Println("==========================")
	}
}
