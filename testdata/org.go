package testdata

import "github.com/xoctopus/sqlx/pkg/types"

// Org 组织
// +genx:model
// @def pk ID
// @def uidx ui_org_id      OrgID;DeletedAt
// @def uidx ui_org_manager OrgID;Manager;DeletedAt
// @def idx  i_name         Name;DeletedAt
// @def idx  i_belonged     Belonged
type Org struct {
	types.AutoIncID

	RelOrg
	OrgData

	types.OperationTime
}

type RelOrg struct {
	// OrgID 组织ID
	OrgID OrgID `db:"f_org_id"`
}

type OrgData struct {
	// Name 组织名称
	Name string `db:"f_name,width=256"`
	// Belonged 组织归属组织ID
	Belonged UserID `db:"f_belongs,default='0'"`
	// Manager 组织管理者ID
	Manager UserID
}

type OrgID uint64
