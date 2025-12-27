package types_test

import (
	"database/sql/driver"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

func TestOperationDatetime(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		ops := types.CreationDatetime{}

		t.Run("AutoMarked", func(t *testing.T) {
			Expect(t, ops.CreatedAt.IsZero(), BeTrue())
			ops.MarkCreatedAt()
			Expect(t, ops.CreatedAt.IsZero(), BeFalse())
		})
		t.Run("UserMarked", func(t *testing.T) {
			ts := sqltime.AsDatetime(time.Now())
			ops.CreatedAt = ts
			ops.MarkCreatedAt()
			Expect(t, ts.Equal(ops.CreatedAt.Unwrap()), BeTrue())
		})
	})

	t.Run("CreationModification", func(t *testing.T) {
		ops := types.OperationDatetime{}

		t.Run("AutoMarked", func(t *testing.T) {
			Expect(t, ops.CreatedAt.IsZero(), BeTrue())
			Expect(t, ops.UpdatedAt.IsZero(), BeTrue())
			ops.MarkCreatedAt()
			Expect(t, ops.CreatedAt.IsZero(), BeFalse())
			Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
		})
		t.Run("UserMarked", func(t *testing.T) {
			ts := sqltime.AsDatetime(time.Now())
			ops.CreatedAt = ts
			ops.UpdatedAt = ts
			Expect(t, ts.Equal(ops.CreatedAt.Unwrap()), BeTrue())
			Expect(t, ts.Equal(ops.UpdatedAt.Unwrap()), BeTrue())
		})
	})

	t.Run("CreationModificationDeletion", func(t *testing.T) {
		ops := types.OperationDatetime{}

		Expect(t, ops.DeletedAt.IsZero(), BeTrue())
		ops.MarkDeletedAt()
		Expect(t, ops.DeletedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.Equal(ops.DeletedAt.Unwrap()), BeTrue())

		col, _, defv := ops.SoftDeletion()
		Expect(t, col, Equal("DeletedAt"))
		Expect(t, defv, Equal[driver.Value](sqltime.DatetimeZero))
	})
}
