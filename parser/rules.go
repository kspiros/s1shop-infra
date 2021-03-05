package parser

import (
	"github.com/kspiros/xlib/parser/actions"
	"github.com/kspiros/xlib/parser/filters"
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
	ValueField    string             `json:"valuefield" bson:"valuefield"`
	CaseSensitive bool               `json:"casesensitive" bson:"casesensitive"`
	Operator      *ConditionOperator `json:"operator,omitempty" bson:"operator,omitempty"`
	Conditions    *[]Condition       `json:"conditions,omitempty" bson:"conditions,omitempty"`
}

type RootCondition struct {
	Operator   ConditionOperator `json:"operator" bson:"operator"`
	Conditions []Condition       `json:"conditions" bson:"conditions"`
}

type Rule struct {
	Name        string           `json:"name" bson:"name"`
	Description string           `json:"description" bson:"description"`
	IsActive    bool             `json:"isactive" bson:"isactive"`
	Rule        RootCondition    `json:"rule" bson:"rule"`
	Actions     []actions.Action `json:"actions" bson:"actions"`
}
