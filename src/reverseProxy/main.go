package reverseProxy

import (
	"../util"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target *url.URL, tr http.RoundTripper) *httputil.ReverseProxy {
	targetRawQuery := target.RawQuery
	targetHost := target.Host

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = targetHost
		req.URL.Path = util.CleanUrlPath(target.Path + "/" + req.URL.Path)
		if len(targetRawQuery) == 0 || len(req.URL.RawQuery) == 0 {
			req.URL.RawQuery = targetRawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetRawQuery + "&" + req.URL.RawQuery
		}

		// host header
		req.Host = targetHost
	}

	return &httputil.ReverseProxy{
		Transport: tr,
		Director:  director,
	}
}
