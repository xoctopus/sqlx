package builder

func MatchColumn(c Col, name string) bool {
	if name == "" {
		return false
	}

	if x := name[0]; x >= 'A' && x <= 'Z' {
		return c.FieldName() == name
	}

	return c.Name() == name
}
