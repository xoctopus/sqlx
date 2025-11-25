package nullable

import (
	"testing"

	. "github.com/xoctopus/x/testx"
)

func BenchmarkNewNullIgnoreScanner(b *testing.B) {
	v := 0
	for i := 0; i < b.N; i++ {
		_ = NewNullIgnoreScanner(&v).Scan(2)
	}
}

func TestNullIgnoreScanner(t *testing.T) {
	t.Run("ScanValue", func(t *testing.T) {
		v := 0
		s := NewNullIgnoreScanner(&v)
		_ = s.Scan(2)
		Expect(t, v, Equal(2))
	})

	t.Run("ScanNil", func(t *testing.T) {
		s := NewNullIgnoreScanner(&Empty{})
		Expect(t, s.Scan(nil), Succeed())

		v := 0
		s = NewNullIgnoreScanner(&v)
		Expect(t, s.Scan(nil), Succeed())
		Expect(t, v, Equal(0))
	})
}
