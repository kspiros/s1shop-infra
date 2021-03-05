package filters

import "strings"

const (
	cStartsWith = "starts_with"
)

type fStartsWith struct {
}

func (t *fStartsWith) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := cvalue.(string); !ok {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}
	return strings.HasPrefix(value.(string), cvalue.(string))
}

func init() {
	RegisterFilter(cStartsWith, &fStartsWith{})
}
