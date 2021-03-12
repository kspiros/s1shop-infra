package actions

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var actionList map[string]reflect.Type

type IAction interface {
	Execute(row *map[string]interface{})
}

type Action struct {
	Action string      `json:"action" bson:"action"`
	Params interface{} `json:"params" bson:"params"`
}

func GetAction(action string) (reflect.Type, bool) {
	c, found := actionList[action]
	return c, found
}

func RegisterAction(action string, r reflect.Type) {
	if actionList == nil {
		actionList = make(map[string]reflect.Type, 0)
	}
	actionList[action] = r
}

func GetActions(list *[]Action) []IAction {
	l := make([]IAction, 0)
	for _, action := range *list {
		a := ActionFactory(&action)
		if a != nil {
			l = append(l, a)
		}
	}
	return l
}

func ActionFactory(action *Action) IAction {
	ty, found := GetAction(action.Action)

	if !found {
		return nil
	}

	value := reflect.New(ty).Interface()
	valueBytes, err := json.Marshal(action.Params)
	if err != nil {
		return nil
	}

	if err = json.Unmarshal(valueBytes, &value); err != nil {
		return nil
	}

	if v, ok := value.(IAction); ok {
		return v
	}
	return nil

}

func (a *Action) Unmarshal(m map[string]interface{}) error {

	actionname, success := m["action"].(string)
	if !success {
		return errors.New("wrong syntax action name missing")
	}

	p := m["params"]
	if p == nil {
		return errors.New("wrong syntax params missing")
	}

	ty, found := GetAction(actionname)
	a.Action = actionname

	if !found {
		return nil //errors.New("Service not found")
	}

	value := reflect.New(ty).Interface()

	valueBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(valueBytes, &value); err != nil {
		return err
	}

	a.Params = value
	return nil
}

func (a *Action) UnmarshalBSON(b []byte) error {
	var j map[string]interface{}
	err := bson.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	return a.Unmarshal(j)
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var j map[string]interface{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	return a.Unmarshal(j)
}
