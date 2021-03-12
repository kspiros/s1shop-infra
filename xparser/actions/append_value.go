package actions

import (
	"fmt"
	"reflect"
)

const (
	acAppendValue string = "append_value"
)

type AppendValue struct {
	Field string `json:"field" bson:"field" validate:"required"`
	Value string `json:"value" bson:"value" validate:"required"`
}

func (a *AppendValue) Execute(row *map[string]interface{}) {
	value := fmt.Sprintf("%v%v", (*row)[a.Field], a.Value)
	value = ResolveVariables(value, row).(string)
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acAppendValue, reflect.TypeOf(AppendValue{}))
}
