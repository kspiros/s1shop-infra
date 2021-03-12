package xparser

type IOperator interface {
	Evaluate(condition *Condition, row *map[string]interface{}, isvalid bool) bool
}

type OrOperator struct {
}

func (o *OrOperator) Evaluate(condition *Condition, row *map[string]interface{}, isvalid bool) bool {
	return isvalid || ExecuteFilter(condition, row)
}

type AndOperator struct {
}

func (o *AndOperator) Evaluate(condition *Condition, row *map[string]interface{}, isvalid bool) bool {
	return isvalid && ExecuteFilter(condition, row)
}

type XorOperator struct {
}

func (o *XorOperator) Evaluate(condition *Condition, row *map[string]interface{}, isvalid bool) bool {
	return isvalid != ExecuteFilter(condition, row)
}

func CreateOperator(operator *ConditionOperator) IOperator {
	if *operator == coAND {
		return &AndOperator{}
	} else if *operator == coOR {
		return &OrOperator{}
	} else if *operator == coOR {
		return &XorOperator{}
	}
	return nil
}
