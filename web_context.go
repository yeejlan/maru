package maru

import (
	"net/http"
)

type WebContext struct {
	request *http.Request
	response http.ResponseWriter
}