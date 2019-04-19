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

func (this *WebContext) Abort(status int, body string) {
	this.W.WriteHeader(status)
	this.W.Write([]byte(body))
}

func (this WebContext) Redirect(url string) {
	this.W.Header().Set("Location", url)
	this.W.WriteHeader(302)
}