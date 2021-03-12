package actions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	acSearchForValue string = "search_for_value"
)

type SearchForValue struct {
	Field       string `json:"field" bson:"field" validate:"required"`
	LookUpField string `json:"lookup_field" bson:"lookup_field" validate:"required"`
	Search      string `json:"search" bson:"search" validate:"required"`
	Replace     string `json:"replace" bson:"replace" validate:"required"`
}

func (a *SearchForValue) Execute(row *map[string]interface{}) {

	sv := ResolveVariables(a.Search, row).(string)
	rv := ResolveVariables(a.Replace, row).(string)

	srcvalues := strings.Split(sv, "\n")
	replacevalues := strings.Split(rv, "\n")

	value := fmt.Sprintf("%v", (*row)[a.LookUpField])

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
	RegisterAction(acSearchForValue, reflect.TypeOf(SearchForValue{}))
}
