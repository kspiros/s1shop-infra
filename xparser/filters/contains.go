package filters

import "strings"

const (
	cContains = "contains"
)

type fContains struct {
}

func (t *fContains) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := cvalue.(string); !ok {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}
	return strings.Contains(value.(string), cvalue.(string))
}

func init() {
	RegisterFilter(cContains, &fContains{})
}
