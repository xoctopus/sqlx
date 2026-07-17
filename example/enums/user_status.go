package enums

// UserStatus 用户状态
type UserStatus int8

const (
	USER_STATUS_UNKNOWN   UserStatus = iota
	USER_STATUS__ACTIVITY            // 活跃
	USER_STATUS__DISABLED            // 关闭
)
