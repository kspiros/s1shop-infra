package filters

import "strings"

const (
	cEndsWith = "ends_with"
)

type fEndsWith struct {
}

func (t *fEndsWith) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := cvalue.(string); !ok {
		return false
	}
	if _, ok := value.(string); !ok {
		return false
	}
	return strings.HasSuffix(value.(string), cvalue.(string))
}

func init() {
	RegisterFilter(cEndsWith, &fEndsWith{})
}
