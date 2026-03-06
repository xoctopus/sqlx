package loggingdriver

import "database/sql/driver"

type DriverOptionApplier func(*options)

type options struct {
	ErrorLevel          func(error) int
	DSNParser           func(string) (string, error)
	ValueHolderReplacer func(string, []driver.NamedValue) (string, []driver.NamedValue)
}

func WithErrorLeveler(level func(error) int) func(*options) {
	return func(o *options) {
		o.ErrorLevel = level
	}
}

func WithDsnParser(parser func(string) (string, error)) func(*options) {
	return func(o *options) {
		o.DSNParser = parser
	}
}

func WithInterpolator(replacer func(string, []driver.NamedValue) (string, []driver.NamedValue)) func(*options) {
	return func(o *options) {
		o.ValueHolderReplacer = replacer
	}
}
