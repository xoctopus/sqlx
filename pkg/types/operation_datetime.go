package types

import (
	"database/sql/driver"
	"time"
)

type CreationDatetime struct {
	// CreatedAt 创建时间
	CreatedAt Datetime `db:"f_created_at,precision=3,default=CURRENT_TIMESTAMP(3)" json:"createdAt"`
}

func (c *CreationDatetime) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt.Time = time.Now()
	}
}

type CreationModificationDatetime struct {
	CreationDatetime
	// UpdatedAt 更新时间
	UpdatedAt Datetime `db:"f_updated_at,precision=3,default=CURRENT_TIMESTAMP(3),onupdate=CURRENT_TIMESTAMP(3)" json:"updatedAt"`
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

type CreationModificationDeletionDatetime struct {
	CreationModificationDatetime
	// DeletedAt 删除时间
	DeletedAt Datetime `db:"f_deleted_at,precision=3,default='0001-01-01 00:00:00'" json:"deletedAt"`
}

func (cud CreationModificationDeletionDatetime) SoftDeletion() (string, driver.Value) {
	return "DeletedAt", DatetimeZero
}

func (cud *CreationModificationDeletionDatetime) MarkDeletedAt() {
	cud.MarkModifiedAt()
	cud.DeletedAt = cud.UpdatedAt
}

type OperationDatetime = CreationModificationDeletionDatetime
