package enums

// OrderStatus 订单状态
type OrderStatus int8

const (
	ORDER_STATUS_UNKNOWN     OrderStatus = iota
	ORDER_STATUS__CREATED                // 已创建
	ORDER_STATUS__PROCESSING             // 处理中
	ORDER_STATUS__FINISHED               // 已完成
	ORDER_STATUS__CANCELED               // 已取消
)
