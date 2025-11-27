package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
	_ "time/tzdata"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"
)

var (
	UTC = time.UTC
	// CST Chinese Standard Timezone
	CST = must.NoErrorV(time.LoadLocation("Asia/Shanghai"))
	// JST Japan Standard Timezone
	JST = must.NoErrorV(time.LoadLocation("Asia/Tokyo"))
	// KST Korea Standard Timezone
	KST = must.NoErrorV(time.LoadLocation("Asia/Seoul"))
	// IST Indochina Standard Timezone (Thailand)
	IST = must.NoErrorV(time.LoadLocation("Asia/Bangkok"))

	TimestampZero     = Timestamp{time.Time{}}
	TimestampUnixZero = Timestamp{time.Unix(0, 0)}
)

type TimestampPrecision int64

const (
	TIMESTAMP_PRECISION__SEC   TimestampPrecision = 1e9 // seconds
	TIMESTAMP_PRECISION__MILLI TimestampPrecision = 1e6 // milliseconds
	TIMESTAMP_PRECISION__MICRO TimestampPrecision = 1e3 // microseconds
)

// gConfig global timestamp config
var gConfig = struct {
	output    *syncx.OnceOverride[string]
	inputs    *syncx.Set[string]
	timezone  *syncx.OnceOverride[*time.Location]
	precision *syncx.OnceOverride[TimestampPrecision]
}{
	output:    syncx.NewOnceOverride(time.DateTime),
	inputs:    syncx.NewSet[string](time.DateTime),
	timezone:  syncx.NewOnceOverride(time.Local),
	precision: syncx.NewOnceOverride(TIMESTAMP_PRECISION__MILLI), // milliseconds as default precision
}

func SetTimestampPrecision(p TimestampPrecision) {
	if p == 1e9 || p == 1e6 || p == 1e3 || p == 1 {
		if p == 1e9 {
			p = 1e6 // the highest precision is time.Microsecond
		}
		gConfig.precision.Set(p)
	}
}

func Unit() time.Duration {
	return time.Duration(gConfig.precision.Value())
}

func GetTimestampPrecision() TimestampPrecision {
	return gConfig.precision.Value()
}

func SetTimestampOutputLayout(layout string) {
	gConfig.output.Set(layout)
}

func GetTimestampOutputLayout() string {
	return gConfig.output.Value()
}

func AddTimestampInputLayouts(layouts ...string) {
	for _, layout := range layouts {
		gConfig.inputs.Store(layout)
	}
}

func GetTimestampInputLayouts() []string {
	return gConfig.inputs.Keys()
}

func SetTimestampTimezone(timezone *time.Location) {
	gConfig.timezone.Set(timezone)
}

func GetTimestampTimezone() *time.Location {
	return gConfig.timezone.Value()
}

func ParseTimestamp(input string) (d Timestamp, err error) {
	for layout := range gConfig.inputs.Range {
		t, e := time.ParseInLocation(layout, input, gConfig.timezone.Value())
		if e == nil {
			return Timestamp{t}, nil
		}
		err = e
	}
	return d, err
}

func ParseTimestampWithLayout(input, layout string) (Timestamp, error) {
	t, err := time.ParseInLocation(layout, input, gConfig.timezone.Value())
	if err != nil {
		return TimestampUnixZero, err
	}
	return Timestamp{t}, nil
}

func AsTimestamp(t time.Time) Timestamp {
	return Timestamp{t}
}

type Timestamp struct {
	time.Time `json:",inline"`
}

func (Timestamp) DBType(driver string) string {
	return "bigint"
}

func (t *Timestamp) Scan(src any) error {
	digital := int64(0)
	switch v := src.(type) {
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("sql.Scan() strfmt.Timestamp from: %#v failed: %s", v, err.Error())
		}
		digital = n
	case int64:
		digital = v
	case nil:
		*t = TimestampZero
		return nil
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Timestamp from: %#v", v)
	}

	if digital < 0 {
		*t = TimestampZero
		return nil
	}
	switch precision := gConfig.precision.Value(); precision {
	case TIMESTAMP_PRECISION__SEC:
		*t = Timestamp{time.Unix(digital, 0)}
	case TIMESTAMP_PRECISION__MILLI:
		*t = Timestamp{time.Unix(digital/1e3, digital%1e3*1e6)}
	default:
		must.BeTrue(precision == TIMESTAMP_PRECISION__MICRO)
		*t = Timestamp{time.Unix(digital/1e6, digital%1e6*1e3)}
	}

	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	unix := t.Int()
	if unix < 0 {
		unix = 0
	}
	return unix, nil
}

func (t Timestamp) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(gConfig.output.Value())
	// return t.In(gConfig.timezone.Value()).Format(gConfig.output.Value())
}

func (t Timestamp) Format(layout string) string {
	return t.In(gConfig.timezone.Value()).Format(layout)
}

func (t Timestamp) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Timestamp) UnmarshalText(data []byte) error {
	s := string(data)
	if len(s) == 0 || s == "0" {
		return nil
	}
	x, err := ParseTimestamp(s)
	if err != nil {
		return err
	}
	*t = x
	return err
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return t.UnmarshalText([]byte(s))
}

func (t Timestamp) Int() int64 {
	switch precision := gConfig.precision.Value(); precision {
	case TIMESTAMP_PRECISION__SEC:
		return t.Time.Unix()
	case TIMESTAMP_PRECISION__MILLI:
		return t.Time.UnixMilli()
	default:
		must.BeTrue(TIMESTAMP_PRECISION__MICRO == precision)
		return t.Time.UnixMicro()
	}
}

func (t Timestamp) IsZero() bool {
	unix := t.Int()
	return unix == 0 || unix == TimestampZero.Int()
}
