package types

import (
	"database/sql/driver"
	"time"
)

type CreationDatetime struct {
	// CreatedAt 创建时间
	CreatedAt Datetime `db:"f_created_at,default='CURRENT_TIMESTAMP(3)'" json:"createdAt"`
}

func (c *CreationDatetime) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = Datetime{Time: time.Now()}
	}
}

type CreationModificationDatetime struct {
	CreationDatetime
	// UpdatedAt 更新时间
	UpdatedAt Datetime `db:"f_updated_at,default='CURRENT_TIMESTAMP(3)',onupdate='CURRENT_TIMESTAMP(3)'" json:"updatedAt"`
}

func (cu *CreationModificationDatetime) MarkCreatedAt() {
	cu.MarkModifiedAt()
	if cu.CreatedAt.IsZero() {
		cu.CreatedAt = cu.UpdatedAt
	}
}

func (cu *CreationModificationDatetime) MarkModifiedAt() {
	if cu.UpdatedAt.IsZero() {
		cu.UpdatedAt = Datetime{Time: time.Now()}
	}
}

type CreationModificationDeletionDatetime struct {
	CreationModificationDatetime
	// DeletedAt 删除时间
	DeletedAt Datetime `db:"f_deleted_at,default='0000-00-00 00:00:00.000'" json:"deletedAt"`
}

func (cud *CreationModificationDeletionDatetime) MarkDeletedAt() {
	cud.MarkModifiedAt()
	cud.DeletedAt = cud.UpdatedAt
}

func (cud CreationModificationDeletionDatetime) SoftDeletion() (string, driver.Value) {
	return "DeletedAt", DatetimeUnixZero
}

type OperationDatetime = CreationModificationDeletionDatetime
