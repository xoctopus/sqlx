package enums

// ProductStatus 产品状态
type ProductStatus int8

const (
	PRODUCT_STATUS_UNKNOWN   ProductStatus = iota
	PRODUCT_STATUS__ON_SALE                // 售卖中
	PRODUCT_STATUS__DISABLED               // 关闭
)
