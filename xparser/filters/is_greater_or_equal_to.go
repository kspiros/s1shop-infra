package filters

import "github.com/kspiros/xlib"

const (
	cIsGreaterOrEqualTo = "is_greater_or_equal_to"
)

type fIsGreaterOrEqualTo struct {
}

func (t *fIsGreaterOrEqualTo) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := xlib.CheckFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := xlib.CheckFloatValue(value)
	if err != nil {
		return false
	}

	return v >= cv
}

func init() {
	RegisterFilter(cIsGreaterOrEqualTo, &fIsGreaterOrEqualTo{})
}
