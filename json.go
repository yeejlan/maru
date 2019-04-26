package maru

import (
	"encoding/json"
	"strings"
)

var(
	JsonEncode = json.Marshal
	JsonDecode = json.Unmarshal
)

//decode json to map[string]interface{}, int as json.Number
func JsonDecodeToMap(data string) (map[string]interface{}, error) {
	d := json.NewDecoder(strings.NewReader(data))
	d.UseNumber()
	var x interface{}
	if err := d.Decode(&x); err != nil {
		return nil, err
	}
	m, ok := x.(map[string]interface{})
	if !ok {
		return nil, NewError("json convert to map error")
	}
	return m, nil
}
