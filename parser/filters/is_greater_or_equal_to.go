package filters

const (
	cIsGreaterOrEqualTo = "is_greater_or_equal_to"
)

type fIsGreaterOrEqualTo struct {
}

func (t *fIsGreaterOrEqualTo) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := checkFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := checkFloatValue(value)
	if err != nil {
		return false
	}

	return v >= cv
}

func init() {
	RegisterFilter(cIsGreaterOrEqualTo, &fIsGreaterOrEqualTo{})
}
