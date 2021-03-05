package parser

import (
	"strings"

	"github.com/kspiros/xlib/parser/actions"
	"github.com/kspiros/xlib/parser/filters"
)

func ExecuteFilter(c *Condition, row *map[string]interface{}) bool {

	var value interface{}
	var condvalue interface{}
	var ok bool

	value, _ = (*row)[c.Field]

	condvalue = c.Value
	if c.ValueField != "" {
		condvalue, _ = (*row)[c.ValueField]
	}

	f, ok := filters.GetFilter(c.Filter)

	if !ok {
		return false
	}

	if !c.CaseSensitive {
		if v, ok := condvalue.(string); ok {
			condvalue = strings.ToLower(v)
		}
		if v, ok := value.(string); ok {
			value = strings.ToLower(v)
		}
	}

	return f.Execute(condvalue, value) != c.Not

}

func executeConditions(operator *ConditionOperator, conditions *[]Condition, row *map[string]interface{}, isvalid *bool) *bool {
	if (conditions == nil) || (operator == nil) {
		return isvalid
	}

	for _, condition := range *conditions {
		isvalid = executeConditions(condition.Operator, condition.Conditions, row, isvalid)
		v := false
		if isvalid == nil {
			v = ExecuteFilter(&condition, row)
		} else {
			v = ExecuteOperator(operator, &condition, row, *isvalid)
		}
		isvalid = &v
	}
	return isvalid
}

func ExecuteOperator(operator *ConditionOperator, condition *Condition, row *map[string]interface{}, isvalid bool) bool {

	o := CreateOperator(operator)
	if o != nil {
		return o.Evaluate(condition, row, isvalid)
	}
	return false

}

func ExecuteActions(list *[]actions.IAction, row *map[string]interface{}) {
	for _, action := range *list {
		action.Execute(row)
	}
}
func executeRule(root *RootCondition, row *map[string]interface{}) bool {
	isvalid := executeConditions(&root.Operator, &root.Conditions, row, nil)
	return (isvalid != nil) && (*isvalid)
}

func ExecuteParser(rules *[]Rule, data *[]map[string]interface{}) {
	for _, rule := range *rules {
		if !rule.IsActive {
			continue
		}

		l := actions.GetActions(&rule.Actions)

		for _, row := range *data {
			if executeRule(&rule.Rule, &row) {
				ExecuteActions(&l, &row)
			}
		}

	}
}
