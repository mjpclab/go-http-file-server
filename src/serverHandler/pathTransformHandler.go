package serverHandler

import (
	"net/http"

	"mjpclab.dev/ghfs/src/util"
)

type pathTransformHandler struct {
	prefixes    []string
	nextHandler http.Handler
}

func stripUrlPrefix(urlPathDir, urlPath, prefix string) string {
	if len(urlPath) == len(prefix) {
		return "/"
	} else if len(prefix) <= 1 {
		return urlPathDir
	} else {
		return urlPathDir[len(prefix):]
	}
}

func (transformer pathTransformHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPath := util.CleanUrlPath(r.URL.Path)
	var urlPathDir string
	if len(urlPath) > 1 && r.URL.Path[len(r.URL.Path)-1] == '/' {
		urlPathDir = urlPath + "/"
	} else {
		urlPathDir = urlPath
	}
	r.URL.RawPath = urlPathDir

	if len(transformer.prefixes) == 0 {
		r.URL.Path = urlPathDir
		transformer.nextHandler.ServeHTTP(w, r)
		return
	}

	for _, prefix := range transformer.prefixes {
		if !util.HasUrlPrefixDir(urlPath, prefix) {
			continue
		}
		r.URL.Path = stripUrlPrefix(urlPathDir, urlPath, prefix)
		transformer.nextHandler.ServeHTTP(w, r)
		return
	}

	defaultHandler.ServeHTTP(w, r)
}

func newPathTransformHandler(prefixes []string, handler http.Handler) http.Handler {
	return pathTransformHandler{prefixes, handler}
}
