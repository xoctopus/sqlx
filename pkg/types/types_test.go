package types

import (
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"
)

func init() {
	SetTimezone(CST)
	SetTimeOutputLayout(DefaultTimeLayout)

	AddTimeInputLayouts(
		DefaultTimeLayout,
		time.DateTime,
		time.RFC3339Nano,
		time.RFC3339,
		time.DateOnly,
	)
}

func TestTimeConfigurations(t *testing.T) {
	ts, err := ParseTimestamp("1988-10-24 07:00:00.123")
	Expect(t, err, Succeed())
	Expect(t, ts.String(), Equal("1988-10-24 07:00:00.123"))

	Expect(t, GetTimeOutputLayout(), Equal(DefaultTimeLayout))
	Expect(t, GetTimezone(), Equal(CST))
	Expect(t, GetTimeInputLayouts(), HaveLen[[]string](5))

	Expect(t, HasSoftDeletion[M1](), BeTrue())
	Expect(t, HasSoftDeletion[M2](), BeTrue())
}

type M1 struct {
	OperationTime
}

func (M1) TableName() string {
	return "m1"
}

type M2 struct {
	OperationDatetime
}

func (M2) TableName() string {
	return "m2"
}
