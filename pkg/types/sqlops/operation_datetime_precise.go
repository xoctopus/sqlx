package sqlops

import (
	"database/sql/driver"
	"time"

	"github.com/xoctopus/sqlx/pkg/types/sqltime"
)

type CreationDatetimeMilli struct {
	// CreatedAt 创建日期时间(毫秒)
	CreatedAt sqltime.Datetime `db:"created_at,precision=3,default=CURRENT_TIMESTAMP(3)" json:"createdAt"`
}

func (c *CreationDatetimeMilli) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt.Time = time.Now()
	}
}

type CreationModificationDatetimePrecise struct {
	CreationDatetimeMilli
	// UpdatedAt 更新日期时间(毫秒)
	UpdatedAt sqltime.Datetime `db:"updated_at,precision=3,default=CURRENT_TIMESTAMP(3),onupdate=CURRENT_TIMESTAMP(3)" json:"updatedAt"`
}

func (cu *CreationModificationDatetimePrecise) MarkCreatedAt() {
	cu.MarkModifiedAt()
	if cu.CreatedAt.IsZero() {
		cu.CreatedAt = cu.UpdatedAt
	}
}

func (cu *CreationModificationDatetimePrecise) MarkModifiedAt() {
	if cu.UpdatedAt.IsZero() {
		cu.UpdatedAt.Time = time.Now()
	}
}

type CreationModificationDeletionDatetimePrecise struct {
	CreationModificationDatetimePrecise
	// DeletedAt 删除日期时间(毫秒)
	DeletedAt sqltime.Datetime `db:"deleted_at,precision=3,default='0001-01-01 00:00:00.000'" json:"deletedAt"`
}

func (cud CreationModificationDeletionDatetimePrecise) SoftDeletion() (string, []string, driver.Value) {
	return "DeletedAt", []string{"UpdatedAt"}, sqltime.DatetimeZero
}

func (cud *CreationModificationDeletionDatetimePrecise) MarkDeletedAt() {
	cud.MarkModifiedAt()
	cud.DeletedAt = cud.UpdatedAt
}

type (
	OperationDatetimePrecise = CreationModificationDeletionDatetimePrecise
)
