package helper_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/helper"
)

type M struct {
	ID int    `db:"id,autoinc"`
	V  string `db:"v,width=32"`
}

func (M) TableName() string {
	return "helper_test"
}

func TestCVsForInsertion(t *testing.T) {
	cols, vals := helper.CVsForInsertion(
		M{ID: 1, V: "1"},
		M{ID: 2, V: "2"},
		M{ID: 3, V: "3"},
	)
	Expect(t, cols.Len(), Equal(1))
	Expect(t, vals, Equal([]any{"1", "2", "3"}))

	cols, vals = helper.CVsForInsertion(
		&M{ID: 1, V: "1"},
		&M{ID: 2, V: "2"},
		&M{ID: 3, V: "3"},
	)
	Expect(t, cols.Len(), Equal(1))
	Expect(t, vals, Equal([]any{"1", "2", "3"}))
}
