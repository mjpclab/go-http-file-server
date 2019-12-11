package serverHandler

import "net/http"

func (h *handler) auth(w http.ResponseWriter, r *http.Request) (success bool) {
	header := w.Header()
	header.Set("WWW-Authenticate", "Basic realm=\""+r.URL.Path+"\"")

	username, password, hasAuthReq := r.BasicAuth()
	if hasAuthReq {
		success = h.users.Auth(username, password)
	}

	if !success {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return
}
