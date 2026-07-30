package types

import "github.com/shopspring/decimal"

// AsDecimal wraps a shopspring decimal.Decimal.
func AsDecimal(v decimal.Decimal) Decimal {
	return Decimal{}
}

// Decimal decimal number
type Decimal struct {
	decimal.Decimal
}

func (d *Decimal) DBType(driver string) string {
	return "decimal"
}

func (Decimal) OpenAPISchemaType() []string {
	return []string{"number"}
}

func (Decimal) OpenAPISchemaFormat() string {
	return "double"
}
