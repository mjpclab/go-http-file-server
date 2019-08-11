package handler

import "net/http"

type handler struct {
	root              string
	defaultFileServer http.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.defaultFileServer.ServeHTTP(w, r)
}

func NewHandler(root string) *handler {
	h := &handler{
		root:              root,
		defaultFileServer: http.FileServer(http.Dir(root)),
	}
	return h
}
