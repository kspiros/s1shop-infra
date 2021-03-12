package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	acMaximumLength string = "maximum_length"
)

type MaximumLength struct {
	Field  string      `json:"field" bson:"field" validate:"required"`
	Length int         `json:"length" bson:"length" validate:"required"`
	Mode   mLengthMode `json:"mode" bson:"mode" validate:"required" `
}

type mLengthMode int

const (
	lTrimWords mLengthMode = iota
	lTrimDirect
)

func (v mLengthMode) String() string {
	var toString = map[mLengthMode]string{
		lTrimWords:  "words",
		lTrimDirect: "direct",
	}
	return toString[v]
}

// MarshalJSON marshals the enum as a quoted json string
func (v mLengthMode) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (v *mLengthMode) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	var toID = map[string]mLengthMode{
		"words":  lTrimWords,
		"direct": lTrimDirect,
	}
	*v = toID[j]
	return nil
}

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func firstNTrimWord(s string, n int) string {
	i := 0
	for j := range s {
		if (i >= n) && (i+1 > len(s) || s[i+1] == ' ') {
			return s[:j]
		}
		i++
	}
	return s
}

func (a *MaximumLength) Execute(row *map[string]interface{}) {

	value := fmt.Sprintf("%v", (*row)[a.Field])
	if a.Mode == lTrimDirect {
		value = firstN(value, a.Length)
	} else if a.Mode == lTrimWords {
		value = firstNTrimWord(value, a.Length)
	}
	(*row)[a.Field] = value
}

func init() {
	RegisterAction(acMaximumLength, reflect.TypeOf(MaximumLength{}))
}
