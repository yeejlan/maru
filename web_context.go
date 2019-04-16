package maru

import (
	"net/http"
)

type WebContext struct {
	R *http.Request
	W http.ResponseWriter
	//current controller
	Controller string
	//current action
	Action string
	Params map[string]string
	Cookies map[string]string
	Sessions map[string]interface{}
}