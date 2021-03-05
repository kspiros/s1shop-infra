package filters

const (
	cIsGreaterThan = "is_greater_than"
)

type fIsGreaterThan struct {
}

func (t *fIsGreaterThan) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := checkFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := checkFloatValue(value)
	if err != nil {
		return false
	}

	return v > cv
}

func init() {
	RegisterFilter(cIsGreaterThan, &fIsGreaterThan{})
}
