package maru

import (
	"fmt"
	"strconv"
)

type StringMap map[string]string


//implement string interface
func (this StringMap) String() string {
	return fmt.Sprintf("%#v", this)
}

//get string value
func (this StringMap) Get(pathStr string, defaultVal ...string) string {
	if val, ok := this[pathStr]; ok {
		return val
	}
	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return ""
}

//get string value, alias for "Get"
func (this StringMap) GetString(pathStr string, defaultVal ...string) string {
	return this.Get(pathStr, defaultVal...)
}

//get int value
func (this StringMap) GetInt(pathStr string, defaultVal ...int) int {
	if val, ok := this[pathStr]; ok {
		if intval, err := strconv.Atoi(val); err == nil {
			return intval
		}
	}

	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return 0
}

//get bool value
func (this StringMap) GetBool(pathStr string, defaultVal ...bool) bool {
	if val, ok := this[pathStr]; ok {
		if boolval, err := strconv.ParseBool(val); err == nil {
			return boolval
		}
	}

	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return false
}