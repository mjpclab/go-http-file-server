package serverHandler

import (
	"errors"
	"net/http"
)

func (h *aliasHandler) needAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"files\"")
}

func (h *aliasHandler) verifyAuth(r *http.Request, needAuth bool) (username string, success bool, err error) {
	user, pass, hasAuthReq := r.BasicAuth()

	if hasAuthReq && h.users.Auth(user, pass) {
		return user, true, nil
	}

	if !needAuth {
		return "", true, nil
	}

	if !hasAuthReq {
		err = errors.New(r.RemoteAddr + " missing auth info")
	} else {
		err = errors.New(r.RemoteAddr + " auth failed")
	}

	return
}

func (h *aliasHandler) authFailed(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Write([]byte("Unauthorized"))
}
