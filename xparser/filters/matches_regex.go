package filters

import (
	"regexp"
)

const (
	cMatchesRegex = "matches_regex"
)

type fMatchesRegex struct {
}

func (t *fMatchesRegex) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := value.(string); !ok {
		return false
	}

	if _, ok := cvalue.(string); !ok {
		return false
	}

	re := regexp.MustCompile(cvalue.(string))
	results := re.FindAllString(value.(string), -1)

	return len(results) > 0
}

func init() {
	RegisterFilter(cMatchesRegex, &fMatchesRegex{})
}
