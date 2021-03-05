package filters

const (
	cIsEqualTo = "is_equal_to"
)

type fIsEqualTo struct {
}

func (t *fIsEqualTo) Execute(cvalue interface{}, value interface{}) bool {
	return (value != nil) && (cvalue != nil) && (cvalue == value)
}

func init() {
	RegisterFilter(cIsEqualTo, &fIsEqualTo{})
}
