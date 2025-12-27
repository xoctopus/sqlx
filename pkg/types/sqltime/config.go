package sqltime

import (
	"time"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"
)

// gConfig global timestamp config
var gConfig = struct {
	output   *syncx.OnceOverride[string]
	inputs   *syncx.Set[string]
	timezone *syncx.OnceOverride[*time.Location]
}{
	output:   syncx.NewOnceOverride(DefaultOutputLayout),
	timezone: syncx.NewOnceOverride(time.Local),
	inputs: syncx.NewSet[string](
		RFC3339,
		RFC3339Milli,
		DateTime,
		DateTimeMilli,
	),
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

	TimestampZero          = Timestamp{time.Time{}}
	TimestampUnixZero      = Timestamp{time.Unix(0, 0)}
	TimestampMilliZero     = TimestampMilli{time.Time{}}
	TimestampMilliUnixZero = TimestampMilli{time.Unix(0, 0)}

	RFC3339       = time.RFC3339
	RFC3339Milli  = "2006-01-02T15:04:05.999Z07:00"
	DateTime      = time.DateTime
	DateTimeMilli = "2006-01-02 15:04:05.000"

	DefaultOutputLayout = RFC3339Milli
)
