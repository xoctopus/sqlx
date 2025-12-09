package types_test

import (
	"database/sql/driver"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestDatetime(t *testing.T) {
	datetime := types.Datetime{}

	t.Run("DBType", func(t *testing.T) {
		Expect(t, datetime.DBType("mysql"), Equal("datetime"))
		Expect(t, datetime.DBType("pg"), Equal("timestamp"))
		ExpectPanic[string](t, func() { datetime.DBType("sqlite") })
	})

	t.Run("Scanner_Valuer", func(t *testing.T) {
		Expect(t, datetime.Scan(nil), Succeed())
		Expect(t, datetime.IsZero(), BeTrue())

		Expect(t, datetime.Scan([]byte("1988-10-24 07:00:00.123")), Succeed())
		Expect(t, datetime.Int(), Equal(int64(593650800123)))

		Expect(t, datetime.Scan([]byte(time.RFC850)), Failed())
		Expect(t, datetime.Int(), Equal(int64(593650800123)))

		Expect(t, datetime.Scan(593650800123), Failed())

		ts := time.Now()
		Expect(t, datetime.Scan(ts), Succeed())
		Expect(t, ts.Equal(datetime.Unwrap()), BeTrue())

		dv, err := datetime.Value()
		Expect(t, err, Succeed())
		Expect(t, dv, Equal[driver.Value](ts))
	})
}
