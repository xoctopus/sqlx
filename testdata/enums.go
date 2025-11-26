package testdata

// Gender
// +genx:enum
// @def storage=text
type Gender int8

const (
	GENDER_UNKNOWN Gender = iota
	// GENDER__MALE
	// @attr name=男
	// @attr text=男
	// @attr short=M
	GENDER__MALE // 男
	// GENDER__FEMALE
	// @attr name=女
	// @attr text=女
	// @attr short=F
	// @ignore GENDER__FEMALE has no more description, use key as its Text
	GENDER__FEMALE
)
