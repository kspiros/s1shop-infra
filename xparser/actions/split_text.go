package actions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	acSplitText string = "split_text"
)

type SplitText struct {
	Field       string `json:"field" bson:"field" validate:"required"`
	LookUpField string `json:"lookup_field" bson:"lookup_field" validate:"required"`
	Split       string `json:"split" bson:"split" validate:"required" `
	From        int    `json:"from" bson:"from" validate:"required" `
	To          int    `json:"to" bson:"to" validate:"required" `
}

func (a *SplitText) Execute(row *map[string]interface{}) {

	lookupvalue := fmt.Sprintf("%v", (*row)[a.LookUpField])

	values := strings.Split(lookupvalue, a.Split)
	from := a.From
	if from < 0 {
		from = len(values) - from + 1
	} else {
		from--
	}

	to := a.To
	if to < 0 {
		to = len(values) - to + 1
	} else {
		to--
	}

	value := ""
	for i := from; i < to; i++ {
		if value != "" {
			value = value + a.Split + values[i]
		} else {
			value = values[i]
		}
	}
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acSplitText, reflect.TypeOf(SplitText{}))
}
