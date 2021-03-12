package filters

import "strings"

const (
	cContainsAnyOf = "contains_any_of"
)

type fContainsAnyOf struct {
}

func (t *fContainsAnyOf) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := cvalue.(string); !ok {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}

	list := strings.Split(cvalue.(string), "\n")
	for _, v := range list {
		if strings.Contains(value.(string), v) {
			return true
		}
	}

	return false
}

func init() {
	RegisterFilter(cContainsAnyOf, &fContainsAnyOf{})
}
