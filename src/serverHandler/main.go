package serverHandler

import (
	"../serverError"
	"net/http"
	"text/template"
)

type handler struct {
	root              string
	template          *template.Template
	defaultFileServer http.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageData := getPageData(h.root, r)

	if pageData.Item != nil && !pageData.Item.IsDir() {
		h.defaultFileServer.ServeHTTP(w, r)
		return
	}

	err := h.template.Execute(w, pageData)
	serverError.CheckError(err)
}

func NewHandler(root string, template *template.Template) *handler {
	h := &handler{
		root:              root,
		template:          template,
		defaultFileServer: http.FileServer(http.Dir(root)),
	}
	return h
}
