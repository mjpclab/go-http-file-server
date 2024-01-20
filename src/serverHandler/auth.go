package serverHandler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const authQueryParam = "auth"

func (h *aliasHandler) needAuth(rawQuery, rawReqPath, reqFsPath string) (needAuth, requestAuth bool) {
	if strings.HasPrefix(rawQuery, authQueryParam) {
		return true, true
	}

	if h.globalAuth {
		return true, false
	}

	return hasUrlOrDirPrefix(h.authUrls, rawReqPath, h.authDirs, reqFsPath), false
}

func (h *aliasHandler) notifyAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"files\"")
}

func (h *aliasHandler) verifyAuth(r *http.Request, needAuth bool) (userid int, username string, err error) {
	user, pass, hasAuthReq := r.BasicAuth()

	if hasAuthReq {
		var success bool
		if userid, username, success = h.users.Auth(user, pass); success {
			return
		}
	}

	if !needAuth {
		return -1, "", nil
	}

	if !hasAuthReq {
		err = errors.New(r.RemoteAddr + " missing auth info")
	} else {
		err = errors.New(r.RemoteAddr + " auth failed")
	}

	return
}

func (h *aliasHandler) redirectWithoutRequestAuth(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	var returnUrl string
	index := strings.Index(r.URL.RawQuery, authQueryParam+"=")
	if index >= 0 {
		returnUrl = r.URL.RawQuery[index+len(authQueryParam)+1:]
		index = strings.LastIndexByte(returnUrl, '&')
		if index >= 0 {
			returnUrl = returnUrl[:index]
		}
		url, err := url.QueryUnescape(returnUrl)
		if err == nil {
			returnUrl = url
		}
	} else {
		returnUrl = r.Header.Get("Referer")
	}
	if len(returnUrl) == 0 {
		returnUrl = session.prefixReqPath + data.Context.QueryString()
	}

	http.Redirect(w, r, returnUrl, http.StatusFound)
}
