package types_test

import (
	"bytes"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sqlx/pkg/types"
)

func TestBoolean(t *testing.T) {
	Expect(t, Boolean(true), Equal(TRUE))
	Expect(t, Boolean(false), Equal(FALSE))

	Expect(t, TRUE.Bool(), BeTrue())
	Expect(t, FALSE.Bool(), BeFalse())

	data, _ := TRUE.MarshalJSON()
	Expect(t, bytes.Equal(data, []byte("true")), BeTrue())
	data, _ = FALSE.MarshalJSON()
	Expect(t, bytes.Equal(data, []byte("false")), BeTrue())
	data, _ = Bool(3).MarshalJSON()
	Expect(t, bytes.Equal(data, []byte("null")), BeTrue())

	b := Boolean(false)
	Expect(t, b.UnmarshalJSON([]byte("true")), Succeed())
	Expect(t, b.Bool(), BeTrue())
	Expect(t, b.UnmarshalJSON([]byte(`"true"`)), Succeed())
	Expect(t, b.Bool(), BeTrue())
	Expect(t, b.UnmarshalJSON([]byte(`"false"`)), Succeed())
	Expect(t, b.Bool(), BeFalse())
	Expect(t, b.UnmarshalJSON([]byte("false")), Succeed())
	Expect(t, b.Bool(), BeFalse())
	Expect(t, b.UnmarshalJSON([]byte("null")), Succeed())
	Expect(t, int(b), Equal(0))
	Expect(t, b.UnmarshalJSON([]byte("invalid")), Failed())
}
