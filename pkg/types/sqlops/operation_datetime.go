package sqlops

import (
	"database/sql/driver"
	"time"

	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

// CreationDatetime holds CreatedAt as a second-precision datetime.
type CreationDatetime struct {
	// CreatedAt 创建日期时间(秒)
	CreatedAt sqltime.Datetime `db:"created_at,default=CURRENT_TIMESTAMP" json:"createdAt"`
}

func (c *CreationDatetime) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt.Time = time.Now()
	}
}

// CreationModificationDatetime adds UpdatedAt to CreationDatetime.
type CreationModificationDatetime struct {
	CreationDatetime
	// UpdatedAt 更新日期时间(秒)
	UpdatedAt sqltime.Datetime `db:"updated_at,default=CURRENT_TIMESTAMP,onupdate=CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (cu *CreationModificationDatetime) MarkCreatedAt() {
	cu.MarkModifiedAt()
	if cu.CreatedAt.IsZero() {
		cu.CreatedAt = cu.UpdatedAt
	}
}

func (cu *CreationModificationDatetime) MarkModifiedAt() {
	if cu.UpdatedAt.IsZero() {
		cu.UpdatedAt.Time = time.Now()
	}
}

// CreationModificationDeletionDatetime adds soft-deletion DeletedAt.
type CreationModificationDeletionDatetime struct {
	CreationModificationDatetime
	// DeletedAt 删除日期时间(秒)
	DeletedAt sqltime.Datetime `db:"deleted_at,default='0001-01-01 00:00:00'" json:"deletedAt"`
}

func (cud CreationModificationDeletionDatetime) SoftDeletion() (string, []string, driver.Value) {
	return "DeletedAt", []string{"UpdatedAt"}, sqltime.DatetimeEpoch
}

func (cud *CreationModificationDeletionDatetime) MarkDeletedAt() {
	cud.MarkModifiedAt()
	cud.DeletedAt = cud.UpdatedAt
}

type (
	// OperationDatetime is an alias of the full create/modify/delete datetime embed.
	OperationDatetime = CreationModificationDeletionDatetime
)
