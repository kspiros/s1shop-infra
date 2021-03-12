package actions

import (
	"regexp"
	"strings"
)

func ResolveVariables(value interface{}, row *map[string]interface{}) interface{} {
	re := regexp.MustCompile("{{(.*?)}}")
	if v, ok := value.(string); ok {
		results := re.FindAllString(v, -1)
		for _, str := range results {
			fld := str[2 : len(str)-2]
			if cell, ok := (*row)[fld].(string); ok {
				value = strings.ReplaceAll(v, str, cell)
			}
		}
	}
	return value
}
