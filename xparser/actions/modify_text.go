package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	strip "github.com/grokify/html-strip-tags-go"
)

const (
	acModifyText string = "maximum_length"
)

type ModifyText struct {
	Field  string      `json:"field" bson:"field" validate:"required"`
	Length int         `json:"length" bson:"length" validate:"required"`
	Mode   mModifyText `json:"mode" bson:"mode" validate:"required" `
}

func (a *ModifyText) Execute(row *map[string]interface{}) {

	value := fmt.Sprintf("%v", (*row)[a.Field])
	switch a.Mode {
	case mCapitalizeFirstChar:
		(*row)[a.Field] = a.capitalizeFirstChar(value)
	case mCapitalizeFirstCharPerWord:
		(*row)[a.Field] = a.capitalizeFirstCharPerWord(value)
	case mCapitalizeFirstCharPerSentence:
		(*row)[a.Field] = a.capitalizeFirstCharPerSentence(value)
	case mLowerCaseAllWords:
		(*row)[a.Field] = a.lowerCaseAllWords(value)
	case mUpperCaseAllWords:
		(*row)[a.Field] = a.upperCaseAllWords(value)
	case mRemoveNonNumericChars:
		(*row)[a.Field] = a.removeNonNumericChars(value)
	case mRemoveDigits:
		(*row)[a.Field] = a.removeDigits(value)
	case mRemoveLineBreaks:
		(*row)[a.Field] = a.removeLineBreaks(value)
	case mRemoveExtraWhiteSpaces:
		(*row)[a.Field] = a.removeExtraWhiteSpaces(value)
	case mRemoveHTMLFromText:
		(*row)[a.Field] = a.removeHTMLFromText(value)
	}

}

func (a *ModifyText) capitalizeFirstChar(value string) interface{} {
	for i, v := range value {
		return string(unicode.ToUpper(v)) + value[i+1:]
	}
	return value
}
func (a *ModifyText) capitalizeFirstCharPerWord(value string) interface{} {
	return strings.Title(strings.ToLower(value))
}
func (a *ModifyText) capitalizeFirstCharPerSentence(value string) interface{} {
	for i, v := range value {
		if i == 0 {
			value = string(unicode.ToUpper(v)) + value[i+1:]
		} else if value[i-1] == '.' {
			value = value[:i-1] + string(unicode.ToUpper(v)) + value[i+1:]
		}
	}
	return value
}
func (a *ModifyText) lowerCaseAllWords(value string) interface{} {
	return strings.ToLower(value)
}
func (a *ModifyText) upperCaseAllWords(value string) interface{} {
	return strings.ToUpper(value)
}
func (a *ModifyText) removeNonNumericChars(value string) interface{} {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		return value
	}
	return reg.ReplaceAllString(value, "")
}
func (a *ModifyText) removeDigits(value string) interface{} {
	reg, err := regexp.Compile("[^0-9.]")
	if err != nil {
		return value
	}
	return reg.ReplaceAllString(value, "")
}
func (a *ModifyText) removeLineBreaks(value string) interface{} {
	return strings.ReplaceAll(value, "\n", "")
}
func (a *ModifyText) removeExtraWhiteSpaces(value string) interface{} {
	return strings.Join(strings.Fields(value), " ")
}
func (a *ModifyText) removeHTMLFromText(value string) interface{} {
	return strip.StripTags(value)
}

type mModifyText int

const (
	mCapitalizeFirstChar mModifyText = iota
	mCapitalizeFirstCharPerWord
	mCapitalizeFirstCharPerSentence
	mLowerCaseAllWords
	mUpperCaseAllWords
	mRemoveNonNumericChars
	mRemoveDigits
	mRemoveLineBreaks
	mRemoveExtraWhiteSpaces
	mRemoveHTMLFromText
)

func (v mModifyText) String() string {
	var toString = map[mModifyText]string{
		mCapitalizeFirstChar:            "capitalize_first_char",
		mCapitalizeFirstCharPerWord:     "capitalize_first_char_per_word",
		mCapitalizeFirstCharPerSentence: "capitalize_first_char_per_sentence",
		mLowerCaseAllWords:              "lowercase_all_words",
		mUpperCaseAllWords:              "uppercase_all_words",
		mRemoveNonNumericChars:          "remove_non_numeric",
		mRemoveDigits:                   "remove_digits",
		mRemoveLineBreaks:               "remove_line_breaks",
		mRemoveExtraWhiteSpaces:         "remove_extra_whitespaces",
		mRemoveHTMLFromText:             "remove_html",
	}
	return toString[v]
}

// MarshalJSON marshals the enum as a quoted json string
func (v mModifyText) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (v *mModifyText) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	var toID = map[string]mModifyText{
		"capitalize_first_char":              mCapitalizeFirstChar,
		"capitalize_first_char_per_word":     mCapitalizeFirstCharPerWord,
		"capitalize_first_char_per_sentence": mCapitalizeFirstCharPerSentence,
		"lowercase_all_words":                mLowerCaseAllWords,
		"uppercase_all_words":                mUpperCaseAllWords,
		"remove_non_numeric":                 mRemoveNonNumericChars,
		"remove_digits":                      mRemoveDigits,
		"remove_line_breaks":                 mRemoveLineBreaks,
		"remove_extra_whitespaces":           mRemoveExtraWhiteSpaces,
		"remove_html":                        mRemoveHTMLFromText,
	}
	*v = toID[j]
	return nil
}

func init() {
	RegisterAction(acModifyText, reflect.TypeOf(ModifyText{}))
}
