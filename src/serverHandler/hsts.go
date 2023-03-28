package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"net/http"
)

func (h *aliasHandler) tryHsts(w http.ResponseWriter, r *http.Request) (needRedirect bool) {
	_, port := util.ExtractHostnamePort(r.Host)

	if len(port) > 0 {
		return
	}

	header := w.Header()
	header.Set("Strict-Transport-Security", "max-age=31536000")

	if r.TLS != nil {
		return
	}

	location := "https://" + r.Host + r.RequestURI
	http.Redirect(w, r, location, getRedirectCode(r))
	return true
}

func (h *aliasHandler) tryToHttps(w http.ResponseWriter, r *http.Request) (needRedirect bool) {
	if r.TLS != nil {
		return
	}

	hostname, _ := util.ExtractHostnamePort(r.Host)

	var targetPort string
	if len(h.toHttpsPort) > 0 && h.toHttpsPort != ":443" {
		targetPort = h.toHttpsPort
	}

	location := "https://" + hostname + targetPort + r.RequestURI
	http.Redirect(w, r, location, getRedirectCode(r))
	return true
}
