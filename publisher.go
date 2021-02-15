package xlib

import (
	"encoding/json"
)

func EncodeMessage(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
