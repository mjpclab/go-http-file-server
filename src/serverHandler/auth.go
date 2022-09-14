package serverHandler

import (
	"errors"
	"net/http"
)

func (h *aliasHandler) needAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"files\"")
}

func (h *aliasHandler) verifyAuth(r *http.Request) (username string, success bool, err error) {
	var password string
	var hasAuthReq bool
	username, password, hasAuthReq = r.BasicAuth()
	if hasAuthReq {
		success = h.users.Auth(username, password)
		if !success {
			err = errors.New(r.RemoteAddr + " auth failed")
		}
	} else {
		err = errors.New(r.RemoteAddr + " missing auth info")
	}

	return
}

func (h *aliasHandler) authFailed(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Write([]byte("Unauthorized"))
}
