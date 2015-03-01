package main

import (
	"io"
	"net/http"
)

// Factory method for helloWorldHandler
func NewDummyHandler() http.Handler {
	return &dummyHandler{}
}

// An http.Handler which says hello world
type dummyHandler struct{}

func (h *dummyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "You are being served with https (hopefully).\n")
}

// http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/SSLDocker.SingleInstance.html
func main() {
	http.Handle("/", NewDummyHandler())
	http.ListenAndServe(":3000", nil)
}
