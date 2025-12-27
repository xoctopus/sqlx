package sqlops

import (
	"database/sql/driver"
	"time"

	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

type CreationTime struct {
	// CreatedAt 创建时间 秒时间戳
	CreatedAt sqltime.Timestamp `db:"created_at,default='0'" json:"createdAt"`
}

func (c *CreationTime) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = sqltime.AsTimestamp(time.Now())
	}
}

type CreationModificationTime struct {
	CreationTime
	// UpdatedAt 更新时间 秒时间戳
	UpdatedAt sqltime.Timestamp `db:"updated_at,default='0'" json:"updatedAt"`
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

type CreationModificationDeletionTime struct {
	CreationModificationTime
	// DeletedAt 删除时间 秒时间戳
	DeletedAt sqltime.Timestamp `db:"deleted_at,default='0'" json:"deletedAt,omitempty"`
}

func (cmd CreationModificationDeletionTime) SoftDeletion() (string, []string, driver.Value) {
	return "DeletedAt", []string{"UpdatedAt"}, int64(0)
}

func (cmd *CreationModificationDeletionTime) MarkDeletedAt() {
	cmd.MarkModifiedAt()
	cmd.DeletedAt = cmd.UpdatedAt
}

type OperationTime = CreationModificationDeletionTime
