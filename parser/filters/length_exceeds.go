package filters

const (
	cLengthExceeds = "length_exceeds"
)

type fLengthExceeds struct {
}

func (t *fLengthExceeds) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := value.(string); !ok {
		return false
	}
	v, err := checkIntValue(cvalue)
	if err != nil {
		return false
	}

	return len(value.(string)) > v
}

func init() {
	RegisterFilter(cLengthExceeds, &fLengthExceeds{})
}
