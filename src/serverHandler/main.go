package serverHandler

import (
	"../serverError"
	"../serverLog"
	"net/http"
	"strings"
	"text/template"
)

type handler struct {
	root      string
	urlPrefix string
	aliases   map[string]string
	uploads   map[string]bool
	template  *template.Template
	logger    *serverLog.Logger
}

func (h *handler) LogRequest(w http.ResponseWriter, r *http.Request) {
	sb := strings.Builder{}

	sb.WriteString(r.RemoteAddr)
	sb.WriteByte(' ')
	sb.WriteString(r.Method)
	sb.WriteByte(' ')
	sb.WriteString(r.RequestURI)

	h.logger.Log(sb.String())
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.LogRequest(w, r)

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
			serverError.LogError(err)
		}()
	}

	if item != nil && !item.IsDir() {
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=0")
	err := h.template.Execute(w, pageData)
	serverError.LogError(err)
}

func NewHandler(
	root string,
	urlPrefix string,
	aliases map[string]string,
	uploads map[string]bool,
	template *template.Template,
	logger *serverLog.Logger,
) *handler {
	h := &handler{
		root:      root,
		urlPrefix: urlPrefix,
		aliases:   aliases,
		uploads:   uploads,
		template:  template,
		logger:    logger,
	}
	return h
}
