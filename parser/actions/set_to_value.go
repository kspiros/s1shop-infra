package actions

import "reflect"

const (
	acSetToValue string = "set_to_value"
)

type SetToValue struct {
	Field string      `json:"field" bson:"field" validate:"required"`
	Value interface{} `json:"value" bson:"value" validate:"required"`
}

func (a *SetToValue) Execute(row *map[string]interface{}) {
	(*row)[a.Field] = a.Value
}

func init() {
	RegisterAction(acSetToValue, reflect.TypeOf(SetToValue{}))
}
