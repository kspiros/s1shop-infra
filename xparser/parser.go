package xparser

import (
	"regexp"
	"strings"

	"github.com/kspiros/xlib/xparser/actions"
	"github.com/kspiros/xlib/xparser/filters"
)

func ResolveVariables(value interface{}, row *map[string]interface{}) interface{} {
	re := regexp.MustCompile("{{(.*?)}}")
	if v, ok := value.(string); ok {
		results := re.FindAllString(v, -1)
		for _, str := range results {
			fld := str[2 : len(str)-2]
			if cell, ok := (*row)[fld].(string); ok {
				value = strings.ReplaceAll(v, str, cell)
			}
		}
	}
	return value
}

func ExecuteFilter(c *Condition, row *map[string]interface{}) bool {

	value, _ := (*row)[c.Field]
	condvalue := ResolveVariables(c.Value, row)

	if !c.CaseSensitive {
		if v, ok := condvalue.(string); ok {
			condvalue = strings.ToLower(v)
		}
		if v, ok := value.(string); ok {
			value = strings.ToLower(v)
		}
	}

	if f, ok := filters.GetFilter(c.Filter); ok {
		return f.Execute(condvalue, value) != c.Not
	}

	return false
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

func calcConditionTotal(conditions *[]Condition, data *[]map[string]interface{}) {
	for i, condition := range *conditions {
		if condition.Conditions != nil {
			calcConditionTotal(condition.Conditions, data)
		}
		if f, ok := filters.FilterSupportsTotals(condition.Filter); ok {
			(*conditions)[i].Value = f.CalcTotals(condition.Field, condition.Value, data)
		}
	}
}

func calcRuleTotals(root *RootCondition, data *[]map[string]interface{}) {
	calcConditionTotal(&root.Conditions, data)
}

func ExecuteParser(rules *[]Rule, data *[]map[string]interface{}) {
	for _, rule := range *rules {
		if !*rule.IsActive {
			continue
		}

		successActionList := actions.GetActions(&rule.OnSuccess)
		failActionList := actions.GetActions(&rule.OnFail)

		calcRuleTotals(rule.Rule, data)

		for _, row := range *data {
			if executeRule(rule.Rule, &row) {
				ExecuteActions(&successActionList, &row)
			} else {
				ExecuteActions(&failActionList, &row)
			}
		}

	}
}
