package serverHandler

import (
	"../util"
	"net/http"
)

func (h *handler) hsts(w http.ResponseWriter, r *http.Request) (needRedirect bool) {
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
	http.Redirect(w, r, location, http.StatusMovedPermanently)
	return true
}

func (h *handler) https(w http.ResponseWriter, r *http.Request) (needRedirect bool) {
	if r.TLS != nil {
		return
	}

	hostname, _ := util.ExtractHostnamePort(r.Host)

	var targetPort string
	if len(h.httpsPort) > 0 && h.httpsPort != ":443" {
		targetPort = h.httpsPort
	}

	location := "https://" + hostname + targetPort + r.RequestURI
	http.Redirect(w, r, location, http.StatusMovedPermanently)
	return true
}
