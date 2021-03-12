package filters

const (
	cIsEmpty = "is_empty"
)

type fIsEmpty struct {
}

func (t *fIsEmpty) Execute(cvalue interface{}, value interface{}) bool {

	_, ok := value.(string)

	return (value == nil) || (ok && (value.(string) == ""))
}

func init() {
	RegisterFilter(cIsEmpty, &fIsEmpty{})
}
