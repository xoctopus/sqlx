package types

import "github.com/shopspring/decimal"

type Decimal struct {
	decimal.Decimal
}

var _ DBValue = (*Decimal)(nil)

func (d *Decimal) DBType(driver string) string {
	return "decimal"
}
