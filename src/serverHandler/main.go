package serverHandler

import (
	"../param"
	"../serverError"
	"../serverLog"
	"net/http"
	"regexp"
	"text/template"
)

type handler struct {
	root      string
	urlPrefix string
	aliases   map[string]string
	uploads   map[string]bool
	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp
	template  *template.Template
	logger    *serverLog.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go h.logRequest(r)

	pageData, notFound, internalError := h.getPageData(r)
	if len(pageData.Errors) > 0 {
		go func() {
			for _, err := range pageData.Errors {
				serverError.LogError(err)
			}
		}()
	}

	if r.Method == "POST" {
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}

	file := pageData.File
	item := pageData.Item

	if file != nil {
		defer func() {
			err := file.Close()
			serverError.LogError(err)
		}()
	}

	if file != nil && item != nil && !item.IsDir() {
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
		return
	}


	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8;")
	header.Set("Cache-Control", "public, max-age=0")
	if internalError {
		w.WriteHeader(http.StatusInternalServerError)
	} else if notFound {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err := h.template.Execute(w, pageData)
	serverError.LogError(err)
}

func NewHandler(
	root string,
	urlPrefix string,
	p *param.Param,
	template *template.Template,
	logger *serverLog.Logger,
) *handler {
	h := &handler{
		root:      root,
		urlPrefix: urlPrefix,
		aliases:   p.Aliases,
		uploads:   p.Uploads,
		shows:     p.Shows,
		showDirs:  p.ShowDirs,
		showFiles: p.ShowFiles,
		hides:     p.Hides,
		hideDirs:  p.HideDirs,
		hideFiles: p.HideFiles,
		template:  template,
		logger:    logger,
	}
	return h
}
