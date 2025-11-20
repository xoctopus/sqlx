package sqltypes_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	expjson "github.com/go-json-experiment/json"
	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sqlx/pkg/sqltypes"
)

func init() {
	SetTimestampTimezone(CST)
	SetTimestampPrecision(TIMESTAMP_PRECISION__DEFAULT)
	SetTimestampOutputLayout(time.DateTime + ".000")

	AddTimestampInputLayouts(time.DateTime + ".000")
	AddTimestampInputLayouts(time.DateTime)
	AddTimestampInputLayouts(time.RFC3339Nano)
	AddTimestampInputLayouts(time.RFC3339)
	AddTimestampInputLayouts(time.DateOnly)
}

func TestTimestamp(t *testing.T) {
	ts, err := ParseTimestamp("1988-10-24")
	Expect(t, err, Succeed())

	data, err := expjson.Marshal(ts)
	Expect(t, err, Succeed())
	Expect(t, string(data), Equal(strconv.Quote(ts.String())))

	data, err = json.Marshal(ts)
	Expect(t, err, Succeed())
	Expect(t, string(data), Equal(strconv.Quote(ts.String())))
	Expect(t, ts.UnmarshalJSON(data), Succeed())

	data2, _ := ts.MarshalText()
	Expect(t, data2, Equal([]byte(ts.String())))

	_, err = ParseTimestamp("1988-10-24 00:00:00")
	Expect(t, err, Succeed())

	// not in output layouts
	_, err = ParseTimestamp(time.RFC850)
	Expect(t, err, Failed())

	_, err = ParseTimestampWithLayout(time.RFC850, time.RFC850)
	Expect(t, err, Succeed())
	_, err = ParseTimestampWithLayout(time.RFC850, time.RFC3339)
	Expect(t, err, Failed())

	ts = TimestampZero
	Expect(t, ts.UnmarshalJSON([]byte(`"0"`)), Succeed())
	Expect(t, ts.IsZero(), BeTrue())
	Expect(t, ts.String(), Equal(""))

	ts = AsTimestamp(TimestampUnixZero.Add(Unit()))
	Expect(t, ts.IsZero(), BeFalse())
	Expect(t, ts.UnmarshalJSON([]byte(`"0"`)), Succeed())
	Expect(t, ts.IsZero(), BeFalse())
	Expect(t, ts.Equal(TimestampUnixZero.Add(Unit())), BeTrue())
}
