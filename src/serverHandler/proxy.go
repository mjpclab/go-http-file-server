package serverHandler

import (
	"../util"
	"net/http"
)

type proxyHandler struct {
	urlPath string
	handler http.Handler
}

type proxyHandlers []*proxyHandler

func getProxyHandler(r *http.Request, proxies proxyHandlers) http.Handler {
	if len(proxies) == 0 {
		return nil
	}

	maxUrlLen := 0
	var proxyHandler http.Handler = nil

	requestUrlPath := r.RequestURI
	for _, proxy := range proxies {
		urlPath := proxy.urlPath
		handler := proxy.handler

		if len(requestUrlPath) < len(urlPath) || !util.HasUrlPrefixDir(requestUrlPath, urlPath) {
			continue
		}

		urlLen := len(urlPath)
		if urlLen > maxUrlLen {
			maxUrlLen = urlLen
			proxyHandler = handler
		}
	}

	return proxyHandler
}

func proxy(w http.ResponseWriter, r *http.Request, proxyHandler http.Handler) {
	w.Header().Set("Cache-Control", "public, max-age=0")
	proxyHandler.ServeHTTP(w, r)
}
