package enums

// ShipmentStatus 物流状态
type ShipmentStatus int8

const (
	SHIPMENT_STATUS_UNKNOWN   ShipmentStatus = iota
	SHIPMENT_STATUS__CREATED                 // 已创建
	SHIPMENT_STATUS__SHIPPING                // 运输中
	SHIPMENT_STATUS__FINISHED                // 已完成
)
