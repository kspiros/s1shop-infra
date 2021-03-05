package filters

import (
	"errors"
	"strconv"
)

func checkIntValue(value interface{}) (int, error) {
	if v, ok := value.(int); ok {
		return v, nil
	} else if v, err := strconv.Atoi(value.(string)); err != nil {
		return v, nil
	}
	return 0, errors.New("no float value")
}

func checkFloatValue(value interface{}) (float64, error) {
	if v, ok := value.(float64); ok {
		return v, nil
	} else if v, err := strconv.ParseFloat(value.(string), 10); err != nil {
		return v, nil
	}
	return 0, errors.New("no float value")
}
