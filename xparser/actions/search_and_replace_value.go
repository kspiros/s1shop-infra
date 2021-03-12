package actions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	acSearchAndReplaceValue string = "search_and_replace_value"
)

type SearchAndReplaceValue struct {
	Field   string `json:"field" bson:"field" validate:"required"`
	Search  string `json:"search" bson:"search" validate:"required"`
	Replace string `json:"replace" bson:"replace" validate:"required" `
}

func (a *SearchAndReplaceValue) Execute(row *map[string]interface{}) {

	sv := ResolveVariables(a.Search, row).(string)
	rv := ResolveVariables(a.Replace, row).(string)

	srcvalues := strings.Split(sv, "\n")
	replacevalues := strings.Split(rv, "\n")

	value := fmt.Sprintf("%v", (*row)[a.Field])

	for i, src := range srcvalues {
		if i < len(replacevalues) {
			value = strings.ReplaceAll(value, src, replacevalues[i])
		} else {
			value = strings.ReplaceAll(value, src, "")
		}
	}
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acSearchAndReplaceValue, reflect.TypeOf(SearchAndReplaceValue{}))
}
