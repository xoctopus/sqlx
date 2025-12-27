package sqlops

import (
	"database/sql/driver"
	"time"

	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

type CreationTimePrecise struct {
	// CreatedAt 创建时间 毫秒时间戳
	CreatedAt sqltime.TimestampMilli `db:"created_at,default='0'" json:"createdAt"`
}

func (c *CreationTimePrecise) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = sqltime.AsTimestampMilli(time.Now())
	}
}

type CreationModificationTimePrecise struct {
	CreationTimePrecise
	// UpdatedAt 更新时间 毫秒时间戳
	UpdatedAt sqltime.TimestampMilli `db:"updated_at,default='0'" json:"updatedAt"`
}

func (cu *CreationModificationTimePrecise) MarkModifiedAt() {
	if cu.UpdatedAt.IsZero() {
		cu.UpdatedAt = sqltime.AsTimestampMilli(time.Now())
	}
}

func (cu *CreationModificationTimePrecise) MarkCreatedAt() {
	cu.MarkModifiedAt()

	if cu.CreatedAt.IsZero() {
		cu.CreatedAt = cu.UpdatedAt
	}
}

type CreationModificationDeletionTimePrecise struct {
	CreationModificationTimePrecise
	// DeletedAt 删除时间 毫秒时间戳
	DeletedAt sqltime.TimestampMilli `db:"deleted_at,default='0'" json:"deletedAt,omitempty"`
}

func (cmd CreationModificationDeletionTimePrecise) SoftDeletion() (string, []string, driver.Value) {
	return "DeletedAt", []string{"UpdatedAt"}, int64(0)
}

func (cmd *CreationModificationDeletionTimePrecise) MarkDeletedAt() {
	cmd.MarkModifiedAt()
	cmd.DeletedAt = cmd.UpdatedAt
}

type OperationTimePrecise = CreationModificationDeletionTimePrecise
