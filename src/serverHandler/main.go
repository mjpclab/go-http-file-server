package serverHandler

import (
	"../serverError"
	"net/http"
	"text/template"
)

type handler struct {
	root     string
	template *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageData := getPageData(h.root, r)
	file := pageData.File
	item := pageData.Item

	if file != nil {
		defer func() {
			err := file.Close()
			serverError.CheckError(err)
		}()
	}

	if item != nil && !item.IsDir() {
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
		return
	}

	err := h.template.Execute(w, pageData)
	serverError.CheckError(err)
}

func NewHandler(root string, template *template.Template) *handler {
	h := &handler{
		root:     root,
		template: template,
	}
	return h
}
