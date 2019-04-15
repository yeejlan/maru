package maru

import (
	"net/http"
	"fmt"
)

type Router struct {

}

func NewRouter() *Router {
	return &Router{}
}

func (this *Router) ServeHTTP(w http.ResponseWriter, request *http.Request){
	fmt.Fprintf(w, "Hello, you've requested: %s\n", request.URL.Path)
}