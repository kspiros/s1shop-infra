package actions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	acLookupAndReplaceValue string = "lookup_and_replace_value"
)

type LookupAndReplaceValue struct {
	Field   string `json:"field" bson:"field" validate:"required"`
	Case    string `json:"case" bson:"case" validate:"required"`
	Replace string `json:"replace" bson:"replace" validate:"required"`
}

func (a *LookupAndReplaceValue) Execute(row *map[string]interface{}) {

	cv := ResolveVariables(a.Case, row).(string)
	rv := ResolveVariables(a.Replace, row).(string)

	casevalues := strings.Split(cv, "\n")
	replacevalues := strings.Split(rv, "\n")

	value := fmt.Sprintf("%v", (*row)[a.Field])

	for i, c := range casevalues {
		if c == value {
			if i < len(replacevalues) {
				value = replacevalues[i]
			} else {
				value = ""
			}
			break
		}
	}
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acLookupAndReplaceValue, reflect.TypeOf(LookupAndReplaceValue{}))
}
