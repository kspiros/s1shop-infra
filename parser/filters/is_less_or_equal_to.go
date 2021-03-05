package filters

const (
	cIsLessOrEqualTo = "is_less_or_equal_to"
)

type fIsLessOrEqualTo struct {
}

func (t *fIsLessOrEqualTo) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := checkFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := checkFloatValue(value)
	if err != nil {
		return false
	}

	return v <= cv
}

func init() {
	RegisterFilter(cIsLessOrEqualTo, &fIsLessOrEqualTo{})
}
