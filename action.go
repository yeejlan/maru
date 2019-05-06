package maru

import (
	"fmt"
	"reflect"
)

//action info
type actionPair struct {
	//controller type, ex. reflect.TypeOf(HomeController{})
	T reflect.Type
	//action, ex. "Index"
	A string 
}

//implement string interface
func (this actionPair) String() string {
	return fmt.Sprintf("ActionPair{T: %T, A: %s}", this.T, this.A)
}

//action storage, for example: 
//ActionMap["home/index"] = ActionPair{T: reflect.TypeOf(HomeController{}), A: "Index",}
var actionMap = make(map[string]actionPair)

//add a action
func AddAction(idx string, T reflect.Type, A string) {
	action := actionPair{
		T: T,
		A: A,
	}
	actionMap[idx] = action
}
