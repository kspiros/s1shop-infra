package filters

const (
	cIsLessThan = "is_less_than"
)

type fIsLessThan struct {
}

func (t *fIsLessThan) Execute(cvalue interface{}, value interface{}) bool {
	cv, err := checkFloatValue(cvalue)
	if err != nil {
		return false
	}

	v, err := checkFloatValue(value)
	if err != nil {
		return false
	}

	return v < cv
}

func init() {
	RegisterFilter(cIsLessThan, &fIsLessThan{})
}
