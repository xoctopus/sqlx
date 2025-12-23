package testdata

import "github.com/xoctopus/sqlx/pkg/types"

// User 用户
// +genx:model
// @attr TableName=t_user
// @def pk ID
// @def u_idx ui_user_id       UserID;DeletedAt
// @def u_idx ui_name          Name;DeletedAt
// @def idx   i_nickname,BTREE Nickname;DeletedAt
// @def idx   i_age            Age
type User struct {
	types.AutoIncID

	RelUser
	RelOrg
	UserData

	types.OperationDatetime
}

type RelUser struct {
	// UserID 用户ID
	UserID UserID `db:"f_user_id"`
}

type UserData struct {
	// Name 用户姓名
	Name string `db:"f_name,width=127"`
	// RealName 真实姓名
	RealName string `db:"f_real_name"`
	// Username 用户姓名
	Username string `db:"f_username,width=255"`
	// Nickname 用户昵称
	Nickname string `db:"f_nick_name,width=127"`
	// Age 年龄
	Age int `db:"f_age"`
	// Gender 性别
	Gender Gender `db:"f_gender"`
	// Asset 资产 decimal(32,4)
	Asset types.Decimal `db:"f_asset,width=32,precision=4"`
}

type UserID uint64
