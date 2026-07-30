package testdata

import "github.com/xoctopus/sqlx/pkg/types"

// User 用户
// +genx:model
// @model TableName=t_user
// @model pk=ID
// @model uidx=ui_user_id;UserID;DeletedAt
// @model uidx=ui_name;Name;DeletedAt
// @model idx=i_nickname,BTREE;Nickname;DeletedAt
// @model idx=i_age;Age
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
	// IsMember 是否为会员
	IsMember bool `db:"f_is_member,default=0"`
	// RealName 真实姓名
	RealName string `db:"f_real_name"`
	// Username 用户姓名
	Username string `db:"f_username,width=255"`
	// Nickname 用户昵称
	Nickname string `db:"f_nick_name,width=128,default=''"`
	// Age 年龄
	Age int `db:"f_age"`
	// Gender 性别
	Gender Gender `db:"f_gender"`
	// Asset 资产
	Asset types.Decimal `db:"f_asset,width=32,precision=4"`
}

type UserID uint64
