package maru

import (
	"fmt"
)

//action info
type actionPair struct {
	//controller instance, ex. HomeController{}
	I interface{}
	//action, ex. "Index"
	A string 
}

//implement string interface
func (this actionPair) String() string {
	return fmt.Sprintf("ActionPair{I: %T, A: %s}", this.I, this.A)
}

//action storage, for example: 
//ActionMap["home/index"] = ActionPair{I: HomeController{}, A: "Index",}
var actionMap = make(map[string]actionPair)

//add a action
func AddAction(idx string, I interface{}, A string) {
	action := actionPair{
		I: I,
		A: A,
	}
	actionMap[idx] = action
}
