package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
	_ "time/tzdata"
)

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

func (t Timestamp) Unwrap() time.Time {
	return t.Time
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
	t.Time = time.Unix(digital/1e3, digital%1e3*1e6)
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
	return t.Time.UnixMilli()
}

func (t Timestamp) IsZero() bool {
	unix := t.Int()
	return unix == 0 || unix == TimestampZero.Int()
}
