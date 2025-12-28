package sqltime

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

var (
	DatetimeZero           = Datetime{Timestamp: TimestampZero}
	DatetimeUnixZero       = Datetime{Timestamp: TimestampUnixZero}
	DatetimeEpoch          = "1970-01-01 00:00:00"
	DatetimeEpochPrecision = "1970-01-01 00:00:00.000"
)

func AsDatetime(t time.Time) Datetime {
	return Datetime{AsTimestamp(t)}
}

type Datetime struct {
	Timestamp
}

func (Datetime) DBType(driver string) string {
	switch strings.ToLower(driver) {
	case "mysql":
		return "datetime"
	case "postgres", "pg":
		return "timestamp"
	default:
		panic("unsupported use Timestamp instead")
	}
}

func (t *Datetime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = DatetimeZero
	case time.Time:
		t.Time = v
	case []byte:
		x, err := ParseTimestamp(string(v))
		if err != nil {
			return err
		}
		t.Time = x.Time
	default:
		return fmt.Errorf("cannot sql.Scan() Datetime from: %#v", v)
	}
	return nil
}

func (t Datetime) Value() (driver.Value, error) {
	return t.Time, nil
}
