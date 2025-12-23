package types

import (
	"database/sql/driver"
	"time"
)

type CreationMarker interface {
	MarkCreatedAt()
}

type ModificationMarker interface {
	MarkModifiedAt()
}

type DeletionMarker interface {
	MarkDeletedAt()
}

type SoftDeletion interface {
	// SoftDeletion returns soft deletion field name and default value
	SoftDeletion() (string, driver.Value)
}

type CreationTime struct {
	// 创建时间
	CreatedAt Timestamp `db:"f_created_at,default='0'" json:"createdAt"`
}

func (c *CreationTime) MarkCreatedAt() {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = Timestamp{time.Now()}
	}
}

type CreationModificationTime struct {
	CreationTime
	// 更新时间
	UpdatedAt Timestamp `db:"f_updated_at,default='0'" json:"updatedAt"`
}

func (cu *CreationModificationTime) MarkModifiedAt() {
	if cu.UpdatedAt.IsZero() {
		cu.UpdatedAt = Timestamp{time.Now()}
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
	// 删除时间
	DeletedAt Timestamp `db:"f_deleted_at,default='0'" json:"deletedAt,omitempty"`
}

func (cmd CreationModificationDeletionTime) SoftDeletion() (string, driver.Value) {
	return "DeletedAt", int64(0)
}

func (cmd *CreationModificationDeletionTime) MarkDeletedAt() {
	cmd.MarkModifiedAt()
	cmd.DeletedAt = cmd.UpdatedAt
}

type OperationTime = CreationModificationDeletionTime
