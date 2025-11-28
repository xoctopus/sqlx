package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

var (
	DatetimeZero     = Datetime{time.Time{}}
	DatetimeUnixZero = Datetime{time.Unix(0, 0)}
)

type Datetime struct {
	time.Time
}

func (Datetime) DBType(driver string) string {
	switch strings.ToLower(driver) {
	case "mysql":
		switch gConfig.precision.Value() {
		case TIMESTAMP_PRECISION__SEC:
			return "datetime"
		case TIMESTAMP_PRECISION__MILLI:
			return "datetime(3)"
		default:
			return "datetime(6)"
		}
	case "postgres", "pg":
		switch gConfig.precision.Value() {
		case TIMESTAMP_PRECISION__SEC:
			return "timestamptz"
		case TIMESTAMP_PRECISION__MILLI:
			return "timestamptz(3)"
		default:
			return "timestamptz(6)"
		}
	default:
		panic("unsupported use Timestamp instead")
	}
}

func (t *Datetime) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*t = DatetimeZero
	case time.Time:
		println("from time.Time", v.Format(time.RFC3339Nano))
		*t = Datetime{v}
	case []byte:
		println("from []byte", string(v))
		x, err := ParseTimestamp(string(v))
		if err != nil {
			return err
		}
		*t = Datetime{x.Time}
	default:
		return fmt.Errorf("cannot sql.Scan() Datetime from: %#v", v)
	}
	return nil
}

func (t Datetime) Value() (driver.Value, error) {
	return t.Time, nil
}
