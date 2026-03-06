package loggingdriver

import (
	"database/sql/driver"
	"strings"
	"testing"
	"time"

	"github.com/xoctopus/x/misc/must"
	. "github.com/xoctopus/x/testx"
)

func TestInterpolate(t *testing.T) {
	// query := "@time = ?"
	for _, c := range []struct {
		value        any
		interpolated string
	}{
		{
			value:        nil,
			interpolated: "NULL",
		},
		{
			value:        must.NoErrorV(time.Parse(time.DateTime+".000000", "1988-10-24 07:00:00.123456")),
			interpolated: "'1988-10-24 07:00:00.123456'",
		},
		{
			value:        time.Time{},
			interpolated: "'0000-00-00 00:00:00.000000'",
		},
		{
			value:        "\x00\n\r\x1a'\"\\ abc",
			interpolated: `'\0\n\r\Z\'\"\\ abc'`,
		},
		{
			value:        []byte("\x00\n\r\x1a'\"\\ abc"),
			interpolated: `E'\0\n\r\Z\'\"\\ abc'`,
		},
		{
			value:        []byte(nil),
			interpolated: "NULL",
		},
		{
			value:        true,
			interpolated: "1",
		},
		{
			value:        false,
			interpolated: "0",
		},
		{
			value:        int64(123),
			interpolated: "123",
		},
		{
			value:        float64(123),
			interpolated: "123",
		},
	} {
		interpolated, err := Interpolate("?c\t", []driver.NamedValue{{Value: c.value}}, time.UTC)
		Expect(t, err, Succeed())
		Expect(t, interpolated, Equal(c.interpolated+"c    "))
	}

	t.Run("UnmatchedArgs", func(t *testing.T) {
		_, err := Interpolate("?,?", []driver.NamedValue{{}}, time.UTC)
		Expect(t, err, IsError(driver.ErrSkip))
	})
	t.Run("ExceedArgLimitation", func(t *testing.T) {
		limit := 65537
		_, err := Interpolate(strings.Repeat("?", limit), make([]driver.NamedValue, limit), time.UTC)
		Expect(t, err, ErrorContains("too many arguments"))
	})
	t.Run("InvalidArgType", func(t *testing.T) {
		_, err := Interpolate("?", []driver.NamedValue{{Value: 1}}, time.UTC)
		Expect(t, err, ErrorContains("unsupported type"))
	})
	t.Run("DefaultInterpolate", func(t *testing.T) {
		DefaultInterpolate("SELECT f_id FROM t_demo WHERE f_name = ?", []driver.NamedValue{{Value: "someone"}})
	})

	printer := NewPrinter("??", []driver.NamedValue{{}})
	Expect(t, printer.String(), HavePrefix("invalid: "))
	printer = NewPrinter("?", []driver.NamedValue{{Value: int64(1)}})
	Expect(t, printer.String(), Equal("1"))
}
