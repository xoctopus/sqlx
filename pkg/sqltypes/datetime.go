package sqltypes

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

var (
	DatetimeZero     = Timestamp{time.Time{}}
	DatetimeUnixZero = Timestamp{time.Unix(0, 0)}
)

type Datetime struct {
	Timestamp
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
		*t = Datetime{TimestampZero}
	case time.Time:
		fmt.Printf("time.Time: %s\n", v.String())
		*t = Datetime{Timestamp{v}}
	case []byte:
		fmt.Printf("bytes: %s\n", string(v))
		x, err := ParseTimestamp(string(v))
		if err != nil {
			return err
		}
		*t = Datetime{x}
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Timestamp from: %#v", v)
	}
	return nil
}

func (t Datetime) Value() (driver.Value, error) {
	return t.Time, nil
}
