package serverHandler

import (
	"../param"
	"../serverErrHandler"
	"../serverLog"
	"../user"
	"html/template"
	"net/http"
	"regexp"
)

type handler struct {
	root      string
	emptyRoot bool
	urlPrefix string

	aliases aliases

	globalUpload bool
	uploadUrls   []string
	uploadDirs   []string

	globalArchive bool
	archiveUrls   []string
	archiveDirs   []string

	globalCors bool
	corsUrls   []string
	corsDirs   []string

	globalAuth bool
	authUrls   []string
	authDirs   []string
	users      user.Users

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp
	template  *template.Template

	logger     *serverLog.Logger
	errHandler *serverErrHandler.ErrHandler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go h.logRequest(r)

	// assert
	if len(r.URL.RawQuery) > 0 {
		queryValues := r.URL.Query()
		assertPath := queryValues.Get("assert")
		if len(assertPath) > 0 {
			h.assert(w, r, assertPath)
			return
		}
	}

	// data
	data := h.getResponseData(r)
	if len(data.errors) > 0 {
		go func() {
			for _, err := range data.errors {
				h.errHandler.LogError(err)
			}
		}()
	}
	file := data.File
	if file != nil {
		defer file.Close()
	}

	if data.NeedAuth && !h.auth(w, r) {
		return
	}

	if data.CanCors {
		h.cors(w, r)
	}

	if data.CanUpload && r.Method == "POST" {
		h.saveUploadFiles(data.handlerReqPath, r)
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// regular flows

	if data.CanArchive && len(r.URL.RawQuery) > 0 {
		switch r.URL.RawQuery {
		case "tar":
			h.tar(w, r, data)
			return
		case "tgz":
			h.tgz(w, r, data)
			return
		case "zip":
			h.zip(w, r, data)
			return
		}
	}

	item := data.Item
	if r.URL.RawQuery == "json" {
		h.json(w, r, data)
	} else if file != nil && item != nil && !item.IsDir() {
		h.content(w, r, data)
	} else {
		h.page(w, r, data)
	}
}

func NewHandler(
	root string,
	emptyRoot bool,
	urlPrefix string,
	p *param.Param,
	users user.Users,
	template *template.Template,
	logger *serverLog.Logger,
	errHandler *serverErrHandler.ErrHandler,
) *handler {
	aliases := aliases{}
	for urlPath, fsPath := range p.Aliases {
		aliases = append(aliases, &alias{urlPath, fsPath})
	}

	h := &handler{
		root:      root,
		emptyRoot: emptyRoot,
		urlPrefix: urlPrefix,

		aliases: aliases,

		globalUpload: p.GlobalUpload,
		uploadUrls:   p.UploadUrls,
		uploadDirs:   p.UploadDirs,

		globalArchive: p.GlobalArchive,
		archiveUrls:   p.ArchiveUrls,
		archiveDirs:   p.ArchiveDirs,

		globalCors: p.GlobalCors,
		corsUrls:   p.CorsUrls,
		corsDirs:   p.CorsDirs,

		globalAuth: p.GlobalAuth,
		authUrls:   p.AuthUrls,
		authDirs:   p.AuthDirs,
		users:      users,

		shows:     p.Shows,
		showDirs:  p.ShowDirs,
		showFiles: p.ShowFiles,
		hides:     p.Hides,
		hideDirs:  p.HideDirs,
		hideFiles: p.HideFiles,
		template:  template,

		logger:     logger,
		errHandler: errHandler,
	}
	return h
}
