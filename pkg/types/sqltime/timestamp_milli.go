package sqltime

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

func ParseTimestampMilli(input string) (d TimestampMilli, err error) {
	for layout := range gConfig.inputs.Range {
		t, e := time.ParseInLocation(layout, input, gConfig.timezone.Value())
		if e == nil {
			return TimestampMilli{t}, nil
		}
		err = e
	}
	return d, err
}

func ParseTimestampMilliWithLayout(input, layout string) (TimestampMilli, error) {
	t, err := time.ParseInLocation(layout, input, gConfig.timezone.Value())
	if err != nil {
		return TimestampMilliZero, err
	}
	return TimestampMilli{t}, nil
}

func AsTimestampMilli(t time.Time) TimestampMilli {
	return TimestampMilli{t.In(gConfig.timezone.Value())}
}

type TimestampMilli struct {
	time.Time `json:",inline"`
}

func (TimestampMilli) DBType(_ string) string {
	return "bigint"
}

func (t TimestampMilli) Unwrap() time.Time {
	return t.Time
}

func (t *TimestampMilli) Scan(src any) error {
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
		*t = TimestampMilliZero
		return nil
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Timestamp from: %#v", v)
	}

	if digital < 0 {
		*t = TimestampMilliZero
		return nil
	}
	t.Time = time.Unix(digital/1e3, digital%1e3*1e6).In(gConfig.timezone.Value())
	return nil
}

func (t TimestampMilli) Value() (driver.Value, error) {
	unix := t.Int()
	if unix < 0 {
		unix = 0
	}
	return unix, nil
}

func (t TimestampMilli) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(gConfig.output.Value())
}

func (t TimestampMilli) Format(layout string) string {
	return t.In(gConfig.timezone.Value()).Format(layout)
}

func (t TimestampMilli) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TimestampMilli) UnmarshalText(data []byte) error {
	s := string(data)
	if len(s) == 0 || s == "0" {
		return nil
	}
	x, err := ParseTimestampMilli(s)
	if err != nil {
		return err
	}
	*t = x
	return err
}

func (t TimestampMilli) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *TimestampMilli) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return t.UnmarshalText([]byte(s))
}

func (t TimestampMilli) Int() int64 {
	return t.Time.UnixMilli()
}

func (t TimestampMilli) IsZero() bool {
	unix := t.Int()
	return unix == 0 || unix == TimestampMilliZero.Int()
}
