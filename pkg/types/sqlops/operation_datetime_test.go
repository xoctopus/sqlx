package sqlops_test

import (
	"database/sql/driver"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types/sqlops"
	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

func TestCreationDatetime(t *testing.T) {
	t.Run("AutoMarked", func(t *testing.T) {
		ops := sqlops.CreationDatetime{}
		Expect(t, ops.CreatedAt.IsZero(), BeTrue())
		ops.MarkCreatedAt()
		Expect(t, ops.CreatedAt.IsZero(), BeFalse())
	})

	t.Run("UserMarked", func(t *testing.T) {
		ops := sqlops.CreationDatetime{}
		ts := sqltime.AsDatetime(time.Now())
		ops.CreatedAt = ts
		ops.MarkCreatedAt()
		Expect(t, ts.Equal(ops.CreatedAt.Unwrap()), BeTrue())
	})
}

func TestCreationModificationDatetime(t *testing.T) {
	t.Run("AutoMarked", func(t *testing.T) {
		ops := sqlops.CreationModificationDatetime{}
		Expect(t, ops.CreatedAt.IsZero(), BeTrue())
		Expect(t, ops.UpdatedAt.IsZero(), BeTrue())

		ops.MarkCreatedAt()
		Expect(t, ops.CreatedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
		Expect(t, ops.CreatedAt.Equal(ops.UpdatedAt.Unwrap()), BeTrue())
	})

	t.Run("UserMarked", func(t *testing.T) {
		ops := sqlops.CreationModificationDatetime{}
		ts := sqltime.AsDatetime(time.Now())
		ops.CreatedAt = ts
		ops.UpdatedAt = ts

		ops.MarkCreatedAt()
		ops.MarkModifiedAt()
		Expect(t, ts.Equal(ops.CreatedAt.Unwrap()), BeTrue())
		Expect(t, ts.Equal(ops.UpdatedAt.Unwrap()), BeTrue())
	})

	t.Run("MarkModifiedAtOnly", func(t *testing.T) {
		ops := sqlops.CreationModificationDatetime{}
		ops.MarkModifiedAt()
		Expect(t, ops.CreatedAt.IsZero(), BeTrue())
		Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
	})
}

func TestCreationModificationDeletionDatetime(t *testing.T) {
	t.Run("MarkDeletedAt", func(t *testing.T) {
		ops := sqlops.OperationDatetime{}
		Expect(t, ops.DeletedAt.IsZero(), BeTrue())

		ops.MarkDeletedAt()
		Expect(t, ops.DeletedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.Equal(ops.DeletedAt.Unwrap()), BeTrue())
	})

	t.Run("SoftDeletion", func(t *testing.T) {
		ops := sqlops.OperationDatetime{}
		col, with, defv := ops.SoftDeletion()
		Expect(t, col, Equal("DeletedAt"))
		Expect(t, with, Equal([]string{"UpdatedAt"}))
		Expect(t, defv, Equal[driver.Value](sqltime.DatetimeEpoch))
	})
}

func TestCreationDatetimeMilli(t *testing.T) {
	t.Run("AutoMarked", func(t *testing.T) {
		ops := sqlops.CreationDatetimeMilli{}
		Expect(t, ops.CreatedAt.IsZero(), BeTrue())
		ops.MarkCreatedAt()
		Expect(t, ops.CreatedAt.IsZero(), BeFalse())
	})

	t.Run("UserMarked", func(t *testing.T) {
		ops := sqlops.CreationDatetimeMilli{}
		ts := sqltime.AsDatetime(time.Now())
		ops.CreatedAt = ts
		ops.MarkCreatedAt()
		Expect(t, ts.Equal(ops.CreatedAt.Unwrap()), BeTrue())
	})
}

func TestCreationModificationDatetimePrecise(t *testing.T) {
	t.Run("AutoMarked", func(t *testing.T) {
		ops := sqlops.CreationModificationDatetimePrecise{}
		Expect(t, ops.CreatedAt.IsZero(), BeTrue())
		Expect(t, ops.UpdatedAt.IsZero(), BeTrue())

		ops.MarkCreatedAt()
		Expect(t, ops.CreatedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
		Expect(t, ops.CreatedAt.Equal(ops.UpdatedAt.Unwrap()), BeTrue())
	})

	t.Run("UserMarked", func(t *testing.T) {
		ops := sqlops.CreationModificationDatetimePrecise{}
		ts := sqltime.AsDatetime(time.Now())
		ops.CreatedAt = ts
		ops.UpdatedAt = ts

		ops.MarkCreatedAt()
		ops.MarkModifiedAt()
		Expect(t, ts.Equal(ops.CreatedAt.Unwrap()), BeTrue())
		Expect(t, ts.Equal(ops.UpdatedAt.Unwrap()), BeTrue())
	})

	t.Run("MarkModifiedAtOnly", func(t *testing.T) {
		ops := sqlops.CreationModificationDatetimePrecise{}
		ops.MarkModifiedAt()
		Expect(t, ops.CreatedAt.IsZero(), BeTrue())
		Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
	})
}

func TestCreationModificationDeletionDatetimePrecise(t *testing.T) {
	t.Run("MarkDeletedAt", func(t *testing.T) {
		ops := sqlops.OperationDatetimePrecise{}
		Expect(t, ops.DeletedAt.IsZero(), BeTrue())

		ops.MarkDeletedAt()
		Expect(t, ops.DeletedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.IsZero(), BeFalse())
		Expect(t, ops.UpdatedAt.Equal(ops.DeletedAt.Unwrap()), BeTrue())
	})

	t.Run("SoftDeletion", func(t *testing.T) {
		ops := sqlops.OperationDatetimePrecise{}
		col, with, defv := ops.SoftDeletion()
		Expect(t, col, Equal("DeletedAt"))
		Expect(t, with, Equal([]string{"UpdatedAt"}))
		Expect(t, defv, Equal[driver.Value](sqltime.DatetimeEpochPrecision))
	})
}
