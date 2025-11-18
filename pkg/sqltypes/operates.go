package sqltypes

import "database/sql/driver"

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
	// SoftDeletion returns soft deletion field and default value
	SoftDeletion() (string, driver.Value)
}

func HasSoftDeletion[M interface{ TableName() string }]() bool {
	_, ok := any(new(M)).(SoftDeletion)
	return ok
}
