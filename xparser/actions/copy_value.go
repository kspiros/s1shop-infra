package actions

import (
	"reflect"
)

const (
	acCopyValue string = "copy_value"
)

type CopyValue struct {
	Field string `json:"field" bson:"field" validate:"required"`
	Value string `json:"value" bson:"value" validate:"required"`
}

func (a *CopyValue) Execute(row *map[string]interface{}) {
	(*row)[a.Field] = (*row)[a.Value]
}

func init() {
	RegisterAction(acCopyValue, reflect.TypeOf(CopyValue{}))
}
