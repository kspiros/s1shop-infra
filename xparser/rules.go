package xparser

import (
	"time"

	"github.com/kspiros/xlib/xparser/actions"
	"github.com/kspiros/xlib/xparser/filters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConditionOperator string

const (
	coAND ConditionOperator = "and"
	coOR                    = "or"
	coXOR                   = "xor"
)

type Condition struct {
	Field         string             `json:"field" bson:"field"`
	Filter        filters.FilterType `json:"filter" bson:"filter"`
	Not           bool               `json:"not" bson:"not"`
	Value         interface{}        `json:"value" bson:"value"`
	CaseSensitive bool               `json:"casesensitive" bson:"casesensitive"`
	Operator      *ConditionOperator `json:"operator,omitempty" bson:"operator,omitempty"`
	Conditions    *[]Condition       `json:"conditions,omitempty" bson:"conditions,omitempty"`
}

type RootCondition struct {
	Operator   ConditionOperator `json:"operator" bson:"operator"`
	Conditions []Condition       `json:"conditions" bson:"conditions"`
}

type Rule struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ChanelID    *primitive.ObjectID `json:"chanelid,omitempty" bson:"chanelid,omitempty"`
	UserID      primitive.ObjectID  `json:"userid" bson:"userid,omitempty"`
	Name        string              `json:"name" bson:"name,omitempty" validate:"required"`
	Description string              `json:"description" bson:"description,omitempty"`
	IsActive    *bool               `json:"isactive" bson:"isactive,omitempty" validate:"required"`
	Rule        *RootCondition      `json:"rule" bson:"rule,omitempty" `
	OnSuccess   []actions.Action    `json:"onsuccess" bson:"onsuccess,omitempty"`
	OnFail      []actions.Action    `json:"onfail" bson:"onfail,omitempty"`
	UpdateDate  time.Time           `json:"updatedate,omitempty" bson:"updatedate,omitempty"`
	InsertDate  time.Time           `json:"insertdate,omitempty" bson:"insertdate,omitempty"`
	Dirty       interface{}         `json:"dirty,omitempty" bson:"dirty,omitempty"`
	Version     *int                `json:"version,omitempty" bson:"version,omitempty"`
}

func (r Rule) GetID() primitive.ObjectID {
	return r.ID
}

func (r Rule) GetUserID() primitive.ObjectID {
	return r.UserID
}

func (r Rule) GetChanelID() primitive.ObjectID {
	return *r.ChanelID
}

func (r Rule) GetVersion() int {
	return *r.Version
}
