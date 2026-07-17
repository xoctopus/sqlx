package enums

// Currency 货币种类
type Currency int8

const (
	CURRENCY_UNKNOWN Currency = iota
	CURRENCY__CNY             // 人民币
	CURRENCY__USD             // 美元
	CURRENCY__JPY             // 日元
)
