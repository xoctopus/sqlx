package mysql

import "strings"

var (
	// MustQuoteBases are the bases that must be quoted.
	MustQuoteBases = map[string]struct{}{
		"CHAR":       {},
		"VARCHAR":    {},
		"TEXT":       {},
		"TINYTEXT":   {},
		"MEDIUMTEXT": {},
		"LONGTEXT":   {},
		"BINARY":     {},
		"VARBINARY":  {},
		"BLOB":       {},
		"TINYBLOB":   {},
		"MEDIUMBLOB": {},
		"LONGBLOB":   {},
		"DATETIME":   {},
	}
	DefaultKeywords = []string{
		"CURRENT_TIMESTAMP",
		"NOW",
		"LOCALTIME",
		"LOCALTIMESTAMP",
		"UTC_TIMESTAMP",
		"NULL",
		"TRUE",
		"FALSE",
	}
	// TextBases can describe width by `information_schema.COLUMNS.character_maximum_length`
	TextBases = map[string]struct{}{
		"CHAR":       {},
		"VARCHAR":    {},
		"TINYTEXT":   {},
		"TEXT":       {},
		"MEDIUMTEXT": {},
		"LONGTEXT":   {},
		// "ENUM"
		// "SET"
	}
	// BinaryBases can describe width by `information_schema.COLUMNS.character_octet_length`
	BinaryBases = map[string]struct{}{
		"BINARY":     {},
		"VARBINARY":  {},
		"TINYBLOB":   {},
		"BLOB":       {},
		"MEDIUMBLOB": {},
		"LONGBLOB":   {},
		// "GEOMETRY"
	}
	// DatetimeBases can describe width by `information_schema.COLUMNS.datetime_precision`
	DatetimeBases = map[string]struct{}{
		"TIME":      {},
		"DATETIME":  {},
		"TIMESTAMP": {},
	}
	// IntegerBases can describe width by `information_schema.COLUMNS.numeric_precision`
	IntegerBases = map[string]struct{}{
		"TINYINT":   {},
		"SMALLINT":  {},
		"MEDIUMINT": {},
		"INT":       {},
		"INTEGER":   {},
		"BIGINT":    {},
		"BIT":       {},
	}
	// FloatBases can describe width by `information_schema.COLUMNS.numeric_precision` and
	// `information_schema.COLUMNS.numeric_precision`
	FloatBases = map[string]struct{}{
		"DECIMAL": {},
		"NUMERIC": {},
		"FLOAT":   {},
		"DOUBLE":  {},
	}
	TextBasesDefaultWidth = map[string]uint64{
		"CHAR":       1,
		"TINYTEXT":   255,
		"TEXT":       65535,
		"MEDIUMTEXT": 16777215,
		"LONGTEXT":   4294967295,
	}
	BinaryBasesDefaultWidth = map[string]uint64{
		"BINARY":     1,
		"TINYBLOB":   255,
		"BLOB":       65535,
		"MEDIUMBLOB": 16777215,
		"LONGBLOB":   4294967295,
	}
	DatetimeBasesDefaultPrecision = map[string]uint64{
		"TIME":      0,
		"DATETIME":  0,
		"TIMESTAMP": 0,
	}
	UnsignedIntegerBasesDefaultWidth = map[string]uint64{
		"TINYINT":   4,
		"SMALLINT":  6,
		"MEDIUMINT": 9,
		"INT":       11,
		"INTEGER":   11,
		"BIGINT":    20,
		"BIT":       1,
	}
	IntegerBasesDefaultWidth = map[string]uint64{
		"TINYINT":   3,
		"SMALLINT":  5,
		"MEDIUMINT": 8,
		"INT":       10,
		"INTEGER":   10,
		"BIGINT":    20,
		"BIT":       1,
	}
	FloatBasesDefaultWidth = map[string]uint64{
		"DECIMAL": 10,
		"NUMERIC": 10,
		"FLOAT":   12,
		"DOUBLE":  22,
	}
	FloatBasesDefaultPrecision = map[string]uint64{
		"DECIMAL": 0,
		"NUMERIC": 0,
	}
)

func DefaultWithWidth(t string, unsigned bool) (string, bool) {
	t = strings.ToUpper(t)
	switch t {
	case "TINYINT":
		if unsigned {
			return t + "(3)", true
		}
		return t + "(4)", true
	case "SMALLINT":
		if unsigned {
			return t + "(5)", true
		}
		return t + "(6)", true
	case "MEDIUMINT":
		if unsigned {
			return t + "(8)", true
		}
		return t + "(9)", true
	case "INT", "INTEGER":
		if unsigned {
			return t + "(10)", true
		}
		return t + "(11)", true
	case "BIGINT":
		if unsigned {
			return t + "(20)", true
		}
		return t + "(19)", true
	case "BIT":
		return t + "(1)", true
	case "TINYTEXT", "TINYBLOB":
		return t + "(255)", true
	case "TEXT", "BLOB":
		return t + "(65536)", true
	case "MEDIUMTEXT", "MEDIUMBLOB":
		return t + "(16777215)", true
	case "LONGTEXT", "LONGBLOB":
		return t + "(4294967295)", true
	case "DATETIME", "TIME", "TIMESTAMP":
		return t + "(0)", true
	case "DECIMAL", "NUMERIC":
		return t + "(10,0)", true
	case "FLOAT":
		return t + "(12,0)", true
	case "DOUBLE":
		return t + "(22,0)", true
	default:
		return t, false
	}
}
