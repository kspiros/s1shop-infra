package xlib

import (
	"errors"
	"sort"
	"strconv"
)

func CheckIntValue(value interface{}) (int, error) {
	if v, ok := value.(int); ok {
		return v, nil
	} else if v, err := strconv.Atoi(value.(string)); err == nil {
		return v, nil
	}
	return 0, errors.New("no float value")
}

func CheckFloatValue(value interface{}) (float64, error) {
	if v, ok := value.(float64); ok {
		return v, nil
	} else if v, err := strconv.ParseFloat(value.(string), 10); err == nil {
		return v, nil
	}
	return 0, errors.New("no float value")
}

func InsertSortedFloat(vs []float64, v float64) []float64 {
	i := sort.SearchFloat64s(vs, v)
	vs = append(vs, 0)
	copy(vs[i+1:], vs[i:])
	vs[i] = v
	return vs
}
