package testdata

// Gender
// +genx:enum
// @enum storage=text
type Gender int8

const (
	GENDER_UNKNOWN Gender = iota
	// GENDER__MALE 男
	// @enum ext.name=男
	// @enum ext.text=男
	// @enum ext.short=M
	GENDER__MALE
	// GENDER__FEMALE 女
	// @enum ext.name=女
	// @enum ext.text=女
	// @enum ext.short=F
	// @enum x=y
	GENDER__FEMALE
)
