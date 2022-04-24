package serverHandler

import "net/http"

func (h *handler) needAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\""+r.URL.Path+"\"")
}

func (h *handler) verifyAuth(r *http.Request) (username string, success bool) {
	var password string
	var hasAuthReq bool
	username, password, hasAuthReq = r.BasicAuth()
	if hasAuthReq {
		success = h.users.Auth(username, password)
	}

	return
}

func (h *handler) authFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
