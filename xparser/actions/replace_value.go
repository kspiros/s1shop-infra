package actions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	acReplaceValue string = "replace_value"
)

type ReplaceValue struct {
	Field   string `json:"field" bson:"field" validate:"required"`
	Search  string `json:"search" bson:"search" validate:"required"`
	Replace string `json:"replace" bson:"replace" validate:"required"`
}

func (a *ReplaceValue) Execute(row *map[string]interface{}) {

	sv := ResolveVariables(a.Search, row).(string)
	rv := ResolveVariables(a.Replace, row).(string)

	value := fmt.Sprintf("%v", (*row)[a.Field])
	value = strings.ReplaceAll(value, sv, rv)
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acReplaceValue, reflect.TypeOf(ReplaceValue{}))
}
