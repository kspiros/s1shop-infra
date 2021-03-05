package filters

import (
	"regexp"
)

const (
	cWordCountExceeds = "word_count_exceeds"
)

type fWordCountExceeds struct {
}

func wordCount(value string) int {
	// Match non-space character sequences.
	re := regexp.MustCompile(`[\S]+`)

	// Find all matches and return count.
	results := re.FindAllString(value, -1)
	return len(results)
}

func (t *fWordCountExceeds) Execute(cvalue interface{}, value interface{}) bool {
	if _, ok := value.(string); !ok {
		return false
	}
	v, err := checkIntValue(cvalue)
	if err != nil {
		return false
	}

	return wordCount(value.(string)) > v
}

func init() {
	RegisterFilter(cWordCountExceeds, &fWordCountExceeds{})
}
