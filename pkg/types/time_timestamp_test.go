package types_test

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sqlx/pkg/types"
)

func TestTimestamp(t *testing.T) {
	ts, err := ParseTimestamp("1988-10-24")
	Expect(t, err, Succeed())

	t.Run("Parse", func(t *testing.T) {
		// not in output layouts
		_, err = ParseTimestamp(time.RFC850)
		Expect(t, err, Failed())

		_, err = ParseTimestampWithLayout(time.RFC850, time.RFC850)
		Expect(t, err, Succeed())
		_, err = ParseTimestampWithLayout(time.RFC850, time.RFC3339)
		Expect(t, err, Failed())

		_, err = ParseTimestamp("1988-10-24 00:00:00")
		Expect(t, err, Succeed())
	})

	t.Run("JSONArshaler", func(t *testing.T) {
		data, err := json.Marshal(ts)
		Expect(t, err, Succeed())
		Expect(t, string(data), Equal(strconv.Quote(ts.String())))

		data, err = json.Marshal(ts)
		Expect(t, err, Succeed())
		Expect(t, string(data), Equal(strconv.Quote(ts.String())))
		Expect(t, ts.UnmarshalJSON(data), Succeed())

		ts = TimestampZero
		Expect(t, ts.UnmarshalJSON([]byte(`"0"`)), Succeed())
		Expect(t, ts.IsZero(), BeTrue())
		Expect(t, ts.String(), Equal(""))

		Expect(t, ts.UnmarshalJSON([]byte(strconv.Quote(time.RFC850))), Failed())
		Expect(t, ts.UnmarshalJSON([]byte(time.RFC850)), Failed())
	})

	t.Run("TextArshaler", func(t *testing.T) {
		ts = AsTimestamp(time.Now())
		data, _ := ts.MarshalText()
		Expect(t, data, Equal([]byte(ts.String())))

		Expect(t, ts.UnmarshalText([]byte(`0`)), Succeed())
		Expect(t, ts.UnmarshalText([]byte(time.RFC850)), Failed())
	})

	t.Run("Scanner", func(t *testing.T) {
		t.Run("Bytes", func(t *testing.T) {
			ts = AsTimestamp(time.Now())
			Expect(t, ts.Scan([]byte("593650800123")), Succeed())
			Expect(t, ts.String(), Equal("1988-10-24 07:00:00.123"))

			Expect(t, ts.Scan([]byte("invalid")), Failed())
		})
		Expect(t, ts.Scan(int64(593650800123)), Succeed())
		Expect(t, ts.String(), Equal("1988-10-24 07:00:00.123"))

		Expect(t, ts.Scan(int64(-1)), Succeed())
		Expect(t, ts.IsZero(), BeTrue())

		Expect(t, ts.Scan(nil), Succeed())
		Expect(t, ts.IsZero(), BeTrue())

		Expect(t, ts.Scan(593650800123), Failed())
	})

	t.Run("Valuer", func(t *testing.T) {
		ts = TimestampZero
		dv, err := ts.Value()
		Expect(t, err, Succeed())
		Expect(t, dv, Equal[driver.Value](int64(0)))

		ts = AsTimestamp(time.Now())
		dv, err = ts.Value()
		Expect(t, err, Succeed())
		Expect(t, dv.(int64) > 0, BeTrue())
	})

	ts = AsTimestamp(TimestampUnixZero.Add(time.Millisecond))
	Expect(t, ts.IsZero(), BeFalse())
	Expect(t, AsTimestamp(TimestampUnixZero.Add(time.Millisecond-1)).IsZero(), BeTrue())
	Expect(t, ts.UnmarshalJSON([]byte(`"0"`)), Succeed())
	Expect(t, ts.IsZero(), BeFalse())
	Expect(t, ts.Equal(TimestampUnixZero.Add(time.Millisecond)), BeTrue())
	Expect(t, ts.DBType("any"), Equal("bigint"))
}
