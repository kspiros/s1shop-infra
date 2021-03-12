package filters

import (
	"github.com/kspiros/xlib"
)

const (
	cIsInHighest = "is_in_highest"
)

type fIsInHighest struct {
}

func (t *fIsInHighest) CalcTotals(field string, cvalue interface{}, data *[]map[string]interface{}) interface{} {
	var values []float64
	v, err := xlib.CheckIntValue(cvalue)
	if err != nil {
		return 0
	}

	for _, row := range *data {
		rowvalue, err := xlib.CheckFloatValue(row[field])
		if err != nil {
			continue
		}
		values = xlib.InsertSortedFloat(values, rowvalue)

	}

	if v > len(values) {
		return 0
	}

	return values[len(values)-v]
}

func (t *fIsInHighest) Execute(cvalue interface{}, value interface{}) bool {
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
	RegisterFilter(cIsInHighest, &fIsInHighest{})
}
