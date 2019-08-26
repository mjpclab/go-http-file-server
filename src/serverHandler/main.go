package serverHandler

import (
	"../serverError"
	"net/http"
	"text/template"
)

type handler struct {
	root      string
	urlPrefix string
	aliases   map[string]string
	uploads   map[string]bool
	template  *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageData := h.getPageData(r)

	if r.Method == "POST" {
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}

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

	w.Header().Set("Cache-Control", "public, max-age=0")
	err := h.template.Execute(w, pageData)
	serverError.CheckError(err)
}

func NewHandler(root, urlPrefix string, aliases map[string]string, uploads map[string]bool, template *template.Template) *handler {
	h := &handler{
		root:      root,
		urlPrefix: urlPrefix,
		aliases:   aliases,
		uploads:   uploads,
		template:  template,
	}
	return h
}
