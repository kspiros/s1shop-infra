package filters

import "github.com/kspiros/xlib"

const (
	cIsLessThan = "is_less_than"
)

type fIsLessThan struct {
}

func (t *fIsLessThan) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := xlib.CheckFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := xlib.CheckFloatValue(value)
	if err != nil {
		return false
	}

	return v < cv
}

func init() {
	RegisterFilter(cIsLessThan, &fIsLessThan{})
}
