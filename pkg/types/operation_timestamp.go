package types

import (
	"database/sql/driver"
	"time"

	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

// CreationTime holds CreatedAt as a second-precision unix timestamp.
type CreationTime struct {
	// CreatedAt 创建时间 秒时间戳
	CreatedAt sqltime.Timestamp `db:"f_created_at,default=0" json:"createdAt"`
}

func (c *CreationTime) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = sqltime.AsTimestamp(time.Now())
	}
}

// CreationModificationTime adds UpdatedAt to CreationTime.
type CreationModificationTime struct {
	CreationTime
	// UpdatedAt 更新时间 秒时间戳
	UpdatedAt sqltime.Timestamp `db:"f_updated_at,default=0" json:"updatedAt"`
}

func (cu *CreationModificationTime) MarkModifiedAt() {
	if cu.UpdatedAt.IsZero() {
		cu.UpdatedAt = sqltime.AsTimestamp(time.Now())
	}
}

func (cu *CreationModificationTime) MarkCreatedAt() {
	cu.MarkModifiedAt()

	if cu.CreatedAt.IsZero() {
		cu.CreatedAt = cu.UpdatedAt
	}
}

// CreationModificationDeletionTime adds soft-deletion DeletedAt.
type CreationModificationDeletionTime struct {
	CreationModificationTime
	// DeletedAt 删除时间 秒时间戳
	DeletedAt sqltime.Timestamp `db:"f_deleted_at,default=0" json:"deletedAt,omitzero"`
}

func (cmd CreationModificationDeletionTime) SoftDeletion() (string, []string, driver.Value) {
	return "DeletedAt", []string{"UpdatedAt"}, int64(0)
}

func (cmd *CreationModificationDeletionTime) MarkDeletedAt() {
	cmd.MarkModifiedAt()
	cmd.DeletedAt = cmd.UpdatedAt
}

// OperationTime is an alias of CreationModificationDeletionTime.
type OperationTime = CreationModificationDeletionTime
