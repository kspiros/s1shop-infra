package filters

import "github.com/kspiros/xlib"

const (
	cIsGreaterThan = "is_greater_than"
)

type fIsGreaterThan struct {
}

func (t *fIsGreaterThan) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := xlib.CheckFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := xlib.CheckFloatValue(value)
	if err != nil {
		return false
	}

	return v > cv
}

func init() {
	RegisterFilter(cIsGreaterThan, &fIsGreaterThan{})
}
