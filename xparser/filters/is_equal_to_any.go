package filters

import "strings"

const (
	cIsEqualToAny = "is_equal_to_any"
)

type fIsEqualToAny struct {
}

func (t *fIsEqualToAny) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := cvalue.(string); !ok {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}

	list := strings.Split(cvalue.(string), "\n")
	for _, v := range list {
		if (value.(string)) == v {
			return true
		}
	}

	return false
}

func init() {
	RegisterFilter(cIsEqualToAny, &fIsEqualToAny{})
}
