package actions

import (
	"reflect"
)

const (
	acCombineValue string = "combine_value"
)

type CombineValue struct {
	Field string `json:"field" bson:"field" validate:"required"`
	Value string `json:"value" bson:"value" validate:"required"`
}

func (a *CombineValue) Execute(row *map[string]interface{}) {
	value := a.Value
	value = ResolveVariables(a.Value, row).(string)
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acCombineValue, reflect.TypeOf(CombineValue{}))
}
