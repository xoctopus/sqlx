package sqltime_test

import (
	"database/sql/driver"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sqlx/pkg/types/sqltime"
)

func Init(t *testing.T) {
	SetTimezone(CST)
	Expect(t, GetTimezone(), Equal(CST))

	Expect(t, GetTimeOutputLayout(), Equal(DefaultOutputLayout))
	SetTimeOutputLayout(DefaultOutputLayout)
	Expect(t, GetTimeOutputLayout(), Equal(DefaultOutputLayout))

	Expect(t, GetTimeInputLayouts(), EquivalentSlice([]string{
		RFC3339,
		RFC3339Milli,
		time.DateTime,
		time.DateTime + ".000",
	}))
	AddTimeInputLayouts(RFC3339)
	AddTimeInputLayouts(RFC3339Milli)
	AddTimeInputLayouts(DefaultOutputLayout)
	Expect(t, GetTimeInputLayouts(), EquivalentSlice([]string{
		RFC3339,
		RFC3339Milli,
		time.DateTime,
		time.DateTime + ".000",
	}))
}

func TestSqlTime_Timestamp(t *testing.T) {
	Init(t)

	t.Run("Parse", func(t *testing.T) {
		seconds := []string{
			"1988-10-24 07:00:00",
			"1988-10-24T07:00:00+08:00",
		}
		millis := []string{
			"1988-10-24 07:00:00.123",
			"1988-10-24T07:00:00.123+08:00",
		}
		fails := []string{time.Layout, time.ANSIC}

		for _, input := range append(seconds, millis...) {
			ts, err := ParseTimestamp(input)
			Expect(t, err, Succeed())
			Expect(t, ts.Unix(), Equal(int64(593650800)))
		}

		for _, input := range millis {
			ts, err := ParseTimestampMilli(input)
			Expect(t, err, Succeed())
			Expect(t, ts.Unix(), Equal(int64(593650800)))
			Expect(t, ts.Int(), Equal(int64(593650800123)))
		}

		for _, input := range fails {
			ts, err := ParseTimestamp(input)
			Expect(t, err, Failed())
			Expect(t, ts, Equal(TimestampZero))

			tp, err := ParseTimestampMilli(input)
			Expect(t, err, Failed())
			Expect(t, tp, Equal(TimestampMilliZero))
		}
	})
	t.Run("ParseWithLayout", func(t *testing.T) {
		ts, err := ParseTimestampWithLayout("24 Oct 88 07:00 +0800", time.RFC822Z)
		Expect(t, err, Succeed())
		Expect(t, ts.Unix(), Equal(int64(593650800)))

		ts, err = ParseTimestampWithLayout("24 Oct 88 07:00 +0800", time.Layout)
		Expect(t, err, Failed())
		Expect(t, ts, Equal(TimestampZero))

		tp, err := ParseTimestampMilliWithLayout("24 Oct 88 07:00 +0800", time.RFC822Z)
		Expect(t, err, Succeed())
		Expect(t, tp.Int(), Equal(int64(593650800000)))

		tp, err = ParseTimestampMilliWithLayout("24 Oct 88 07:00 +0800", time.Layout)
		Expect(t, err, Failed())
		Expect(t, tp, Equal(TimestampMilliZero))
	})

	tt := time.Unix(593650800, 123000000).In(GetTimezone())
	ts := AsTimestamp(tt)
	tp := AsTimestampMilli(tt)

	t.Run("As", func(t *testing.T) {
		Expect(t, AsTimestamp(tt).Int(), Equal(int64(593650800)))
		Expect(t, AsTimestampMilli(tt).Int(), Equal(int64(593650800123)))
	})

	t.Run("DBType", func(t *testing.T) {
		Expect(t, ts.DBType("any"), Equal("bigint"))
		Expect(t, tp.DBType("any"), Equal("bigint"))
	})

	t.Run("Unwrap", func(t *testing.T) {
		Expect(t, ts.Unwrap(), Equal(tt))
		Expect(t, tp.Unwrap(), Equal(tt))
	})

	t.Run("ScannerValuer", func(t *testing.T) {
		t.Run("Timestamp", func(t *testing.T) {
			inputs := []any{
				int64(593650800),
				[]byte(`593650800`),
			}
			x := Timestamp{}
			Expect(t, x.Scan(nil), Succeed())
			Expect(t, x, Equal(TimestampZero))
			Expect(t, x.Scan([]byte("invalid")), Failed())
			Expect(t, x.Scan("invalid"), Failed())
			Expect(t, x.Scan(int64(-1)), Succeed())
			Expect(t, x, Equal(TimestampZero))

			for _, input := range inputs {
				expect, _ := ParseTimestamp("1988-10-24 07:00:00")
				err := x.Scan(input)
				Expect(t, err, Succeed())
				Expect(t, x, Equal(expect))
				Expect(t, x.Int(), Equal(int64(593650800)))

				dv, err := x.Value()
				Expect(t, err, Succeed())
				Expect(t, dv, Equal[driver.Value](int64(593650800)))
			}

			dv, err := TimestampZero.Value()
			Expect(t, err, Succeed())
			Expect(t, dv, Equal[driver.Value](int64(0)))
		})
		t.Run("TimestampMilli", func(t *testing.T) {
			inputs := []any{
				int64(593650800123),
				[]byte(`593650800123`),
			}
			y := TimestampMilli{}
			Expect(t, y.Scan(nil), Succeed())
			Expect(t, y, Equal(TimestampMilliZero))
			Expect(t, y.Scan([]byte("invalid")), Failed())
			Expect(t, y.Scan("invalid"), Failed())
			Expect(t, y.Scan(int64(-1)), Succeed())
			Expect(t, y, Equal(TimestampMilliZero))

			for _, input := range inputs {
				expect, _ := ParseTimestampMilli("1988-10-24 07:00:00.123")
				err := y.Scan(input)
				Expect(t, err, Succeed())
				Expect(t, y, Equal(expect))
				Expect(t, y.Int(), Equal(int64(593650800123)))

				dv, err := y.Value()
				Expect(t, err, Succeed())
				Expect(t, dv, Equal[driver.Value](int64(593650800123)))
			}

			dv, err := TimestampMilliZero.Value()
			Expect(t, err, Succeed())
			Expect(t, dv, Equal[driver.Value](int64(0)))
		})
	})

	t.Run("Arshaler", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			Expect(t, ts.String(), Equal("1988-10-24T07:00:00.123+08:00"))
			Expect(t, tp.String(), Equal("1988-10-24T07:00:00.123+08:00"))
		})

		text := []byte("1988-10-24T07:00:00.123+08:00")
		quoted := []byte(`"1988-10-24T07:00:00.123+08:00"`)

		t.Run("MarshalText", func(t *testing.T) {
			data, err := ts.MarshalText()
			Expect(t, err, Succeed())
			Expect(t, data, Equal(text))

			data, err = tp.MarshalText()
			Expect(t, err, Succeed())
			Expect(t, data, Equal(text))
		})

		t.Run("MarshalJSON", func(t *testing.T) {
			data, err := ts.MarshalJSON()
			Expect(t, err, Succeed())
			Expect(t, data, Equal(quoted))

			data, err = tp.MarshalJSON()
			Expect(t, err, Succeed())
			Expect(t, data, Equal(quoted))
		})

		t.Run("UnmarshalText", func(t *testing.T) {
			x1 := Timestamp{}
			Expect(t, x1.UnmarshalText(text), Succeed())
			Expect(t, x1, Equal(ts))
			Expect(t, x1.UnmarshalText([]byte{}), Succeed())
			Expect(t, x1.UnmarshalText([]byte("0")), Succeed())
			Expect(t, x1.UnmarshalText([]byte("invalid")), Failed())

			x2 := TimestampMilli{}
			Expect(t, x2.UnmarshalText(text), Succeed())
			Expect(t, x2, Equal(tp))
			Expect(t, x2.UnmarshalText([]byte{}), Succeed())
			Expect(t, x2.UnmarshalText([]byte("0")), Succeed())
			Expect(t, x2.UnmarshalText([]byte("invalid")), Failed())
		})

		t.Run("UnmarshalJSON", func(t *testing.T) {
			x1 := Timestamp{}
			Expect(t, x1.UnmarshalJSON(quoted), Succeed())
			Expect(t, x1, Equal(ts))

			x2 := TimestampMilli{}
			Expect(t, x2.UnmarshalJSON(quoted), Succeed())
			Expect(t, x2, Equal(tp))

			invalid := []byte(`"123`)
			Expect(t, x1.UnmarshalJSON(invalid), Failed())
			Expect(t, x2.UnmarshalJSON(invalid), Failed())
		})

		t.Run("Zero", func(t *testing.T) {})
		Expect(t, TimestampZero.String(), Equal(""))
		Expect(t, TimestampMilliZero.String(), Equal(""))

	})
}

func TestSqlTime_Datetime(t *testing.T) {
	Init(t)

	tt := time.Unix(593650800, 123000000).In(GetTimezone())
	datetime := Datetime{}
	t.Run("DBType", func(t *testing.T) {
		Expect(t, datetime.DBType("mysql"), Equal("datetime"))
		Expect(t, datetime.DBType("postgres"), Equal("timestamp"))
		ExpectPanic[string](t, func() { datetime.DBType("sqlite") })
	})

	t.Run("ScannerValuer", func(t *testing.T) {
		inputs := []any{
			tt,
			[]byte("1988-10-24 07:00:00.123"),
		}
		for _, input := range inputs {
			Expect(t, datetime.Scan(input), Succeed())
			Expect(t, datetime, Equal(AsDatetime(tt)))

			dv, err := datetime.Value()
			Expect(t, err, Succeed())
			Expect(t, dv, Equal[driver.Value](tt))
		}

		t.Run("ScanFromNil", func(t *testing.T) {
			Expect(t, datetime.Scan(nil), Succeed())
			Expect(t, datetime, Equal(DatetimeZero))
		})

		t.Run("ScanFromInvalid", func(t *testing.T) {
			Expect(t, datetime.Scan([]byte("invalid_layout")), Failed())
			Expect(t, datetime.Scan("invalid_type"), Failed())
		})
	})
}
