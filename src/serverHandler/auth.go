package serverHandler

import (
	"../shimgo"
	"net/http"
)

func (h *handler) auth(w http.ResponseWriter, r *http.Request) (success bool) {
	header := w.Header()
	header.Set("WWW-Authenticate", "Basic realm=\""+r.URL.Path+"\"")

	username, password, hasAuthReq := shimgo.Net_Http_BasicAuth(r)
	if hasAuthReq {
		success = h.users.Auth(username, password)
	}

	if !success {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}
