// Package v2 for testing
// +genx:doc
package v2

import (
	"github.com/xoctopus/sqlx/pkg/types"
	"github.com/xoctopus/sqlx/testdata"
)

// User 用户表V2
// +genx:model
// @model TableName=t_user
// @model pk=ID
// @model uidx=ui_user_id;UserID
// @model uidx=ui_name;RealName
// @model idx=i_gender;Gender
type User struct {
	types.AutoIncID

	testdata.RelUser
	testdata.RelOrg
	UserDataV2

	types.CreationModificationDatetime
}

type UserDataV2 struct {
	// Name deprecated to RealName
	Name string `db:"f_name,deprecated=f_real_name"`
	// RealName 用户真实姓名
	RealName string `db:"f_real_name,width=255,default=''"` // width changed from 127 to 255
	// Age 年龄
	Age int8 `db:"f_age,default=0"` // datatype changed from int to int8
	// Username deprecated
	Username string `db:"f_username,deprecated"`
	// Nickname 用户昵称 deprecated has indexed before
	Nickname string `db:"f_nick_name,deprecated"`
	// Gender 性别
	Gender testdata.Gender `db:"f_gender"` // no change
	// Desc 描述
	Desc string `db:"f_desc"` // new added
}
