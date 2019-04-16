package maru

import (
	"net/http"
	"fmt"
)

//action info
type ActionPair struct {
	//controller instance, ex. HomeController{}
	I interface{}
	//action, ex. "Index"
	A string 
}

//implement string interface
func (this *ActionPair) String() string {
	return fmt.Sprintf("ActionPair{I: %T, A: %s}", this.I, this.A)
}

//action storage, for example: 
//ActionMap["home/index"] = ActionPair{I: HomeController{}, A: "Index",}
var ActionMap = make(map[string]*ActionPair)

//add a action
func AddAction(idx string, I interface{}, A string) {
	action := &ActionPair{
		I: I,
		A: A,
	}
	ActionMap[idx] = action
}

type Router struct {

}

func NewRouter() *Router {
	return &Router{}
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}