package filters

import "github.com/kspiros/xlib"

const (
	cIsLessOrEqualTo = "is_less_or_equal_to"
)

type fIsLessOrEqualTo struct {
}

func (t *fIsLessOrEqualTo) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := xlib.CheckFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := xlib.CheckFloatValue(value)
	if err != nil {
		return false
	}

	return v <= cv
}

func init() {
	RegisterFilter(cIsLessOrEqualTo, &fIsLessOrEqualTo{})
}
