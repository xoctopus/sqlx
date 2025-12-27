package types

import (
	"testing"

	. "github.com/xoctopus/x/testx"
)

func TestTimeConfigurations(t *testing.T) {
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
