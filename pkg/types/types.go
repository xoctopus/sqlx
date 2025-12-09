package types

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/internal"
)

// DBValue can convert between rdb value and go value with description of rdb datatype
type DBValue interface {
	driver.Valuer
	sql.Scanner
	DBType(driver string) string
}

type DBTypeAdapter interface {
	WithDBType(driver string)
}

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
	DefaultTimeLayout = "2006-01-02 15:04:05.000"
)

// type TimePrecision time.Duration
//
// const (
// 	TIME_PRECISION__SECOND = TimePrecision(time.Second)
// 	TIME_PRECISION__MILLI  = TimePrecision(time.Millisecond)
// 	TIME_PRECISION__MICRO  = TimePrecision(time.Microsecond)
// )

// func SetTimePrecision(p TimePrecision) {
// 	gConfig.precision.Set(p)
// }
//
// func GetTimePrecision() TimePrecision {
// 	return gConfig.precision.Value()
// }

// gConfig global timestamp config
var gConfig = struct {
	output   *syncx.OnceOverride[string]
	inputs   *syncx.Set[string]
	timezone *syncx.OnceOverride[*time.Location]
	// precision *syncx.OnceOverride[TimePrecision]
}{
	output:   syncx.NewOnceOverride(time.DateTime),
	inputs:   syncx.NewSet[string](time.DateTime),
	timezone: syncx.NewOnceOverride(time.Local),
	// precision: syncx.NewOnceOverride(TIME_PRECISION__MILLI),
}

func SetTimeOutputLayout(layout string) {
	gConfig.output.Set(layout)
}

func GetTimeOutputLayout() string {
	return gConfig.output.Value()
}

func AddTimeInputLayouts(layouts ...string) {
	for _, layout := range layouts {
		gConfig.inputs.Store(layout)
	}
}

func GetTimeInputLayouts() []string {
	return gConfig.inputs.Keys()
}

func SetTimezone(timezone *time.Location) {
	gConfig.timezone.Set(timezone)
}

func GetTimezone() *time.Location {
	return gConfig.timezone.Value()
}

func HasSoftDeletion[M internal.Model]() bool {
	_, ok := any(new(M)).(SoftDeletion)
	return ok
}
