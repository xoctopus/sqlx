package loggingdriver

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Interpolator(q string, args []driver.NamedValue) fmt.Stringer {
	return &QueryPrinter{q: q, args: args}
}

type QueryPrinter struct {
	q    string
	args []driver.NamedValue
}

func (p *QueryPrinter) String() string {
	s, err := Interpolate(p.q, p.args, time.Local)
	if err != nil {
		return "invalid: " + err.Error()
	}
	return s
}

const (
	digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
)

func Interpolate(q string, args []driver.NamedValue, loc *time.Location) (string, error) {
	// TODO
	// placeholder is literal character: SELECT x as '?' FROM t WHERE name like '%?%'
	// driver.NamedValue unfolded: {Name: @paramN Value: func(n int) string{ return strings.Repeat(".?", n)[1:] }()
	if strings.Count(q, "?") != len(args) {
		return "", driver.ErrSkip
	}

	if len(args) > 65536 {
		return "", fmt.Errorf("too many arguments: %d, query: %s", len(args), q)
	}

	idx := 0
	buf := bytes.NewBuffer(nil)

	for i := range q {
		switch c := q[i]; c {
		case '?':
			arg := args[idx].Value
			idx++
			switch v := arg.(type) {
			case nil:
				buf.WriteString("NULL")
			case []byte:
				if v == nil {
					buf.WriteString("NULL")
				} else {
					buf.WriteString("E'")
					escape(buf, v)
					buf.WriteByte('\'')
				}
			case string:
				buf.WriteByte('\'')
				escape(buf, []byte(v))
				buf.WriteByte('\'')
			case int64:
				buf.WriteString(strconv.FormatInt(v, 10))
			case float64:
				buf.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
			case bool:
				if v {
					buf.WriteByte('1')
				} else {
					buf.WriteByte('0')
				}
			case time.Time:
				if v.IsZero() {
					buf.WriteString("'0000-00-00 00:00:00.000000'")
					continue
				}
				v = v.In(loc).Add(500) // round microseconds => YYYY-MM-DD HH:MM::SS.microseconds
				yr := v.Year()
				yh := yr / 100
				yl := yr % 100
				mo := v.Month()
				dd := v.Day()
				hh := v.Hour()
				mm := v.Minute()
				ss := v.Second()
				ms := v.Nanosecond() / 1e3
				ms65 := ms / 10000
				ms43 := ms / 100 % 100
				ms21 := ms % 100

				buf.Write([]byte{
					'\'',
					digits10[yh], digits01[yh],
					digits10[yl], digits01[yl],
					'-',
					digits10[mo], digits01[mo],
					'-',
					digits10[dd], digits01[dd],
					' ',
					digits10[hh], digits01[hh],
					':',
					digits10[mm], digits01[mm],
					':',
					digits10[ss], digits01[ss],
					'.',
					digits10[ms65], digits01[ms65],
					digits10[ms43], digits01[ms43],
					digits10[ms21], digits01[ms21],
					'\'',
				})
			default:
				return "", fmt.Errorf("unsupported type: %T: %v", v, v)
			}
		case '\t':
			buf.WriteString("    ")
		default:
			buf.WriteByte(c)
		}
	}

	// if idx != len(args) {
	// 	return "", driver.ErrSkip
	// }
	return buf.String(), nil
}

var escapes = map[byte]string{
	'\x00': `\0`,
	'\n':   `\n`,
	'\r':   `\r`,
	'\x1a': `\Z`,
	'\'':   `\'`,
	'"':    `\"`,
	'\\':   `\\`,
}

func escape(buf *bytes.Buffer, v []byte) {
	for _, c := range v {
		switch c {
		case '\x00', '\n', '\r', '\x1a', '\'', '"', '\\':
			buf.WriteString(escapes[c])
		default:
			buf.WriteByte(c)
		}
	}
}
