package types

import "github.com/shopspring/decimal"

func AsDecimal(v decimal.Decimal) Decimal {
	return Decimal{}
}

type Decimal struct {
	decimal.Decimal
}

func (d *Decimal) DBType(driver string) string {
	return "decimal"
}
