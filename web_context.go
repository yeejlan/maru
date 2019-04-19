package maru

import (
	"net/http"
)

//web context
type WebContext struct {
	Req *http.Request
	W http.ResponseWriter
	//current controller
	Controller string
	//current action
	Action string
	Param StringMap
	Cookie StringMap
	Session map[string]interface{}
}

func newWebContext(w http.ResponseWriter, req *http.Request) *WebContext {
	return &WebContext{
		Req: req,
		W: w,
	}
}