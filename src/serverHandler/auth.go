package serverHandler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const authQueryParam = "auth"

func (h *aliasHandler) needAuth(rawQuery, vhostReqPath, reqFsPath string) (needAuth, requestAuth bool) {
	if strings.HasPrefix(rawQuery, authQueryParam) {
		return true, true
	}

	if h.globalAuth {
		return true, false
	}

	if hasUrlOrDirPrefix(h.authUrls, vhostReqPath, h.authDirs, reqFsPath) {
		return true, false
	}

	if matchPath, _ := hasUrlOrDirPrefixUsers(h.authUrlsUsers, vhostReqPath, h.authDirsUsers, reqFsPath, -1); matchPath {
		return true, false
	}

	return false, false
}

func (h *aliasHandler) notifyAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"files\"")
}

func (h *aliasHandler) verifyAuth(r *http.Request, needAuth bool, vhostReqPath, reqFsPath string) (userid int, username string, err error) {
	user, pass, hasAuthReq := r.BasicAuth()

	if hasAuthReq {
		var success bool
		userid, username, success = h.users.Auth(user, pass)
		if success && userid >= 0 && (len(h.authUrlsUsers) > 0 || len(h.authDirsUsers) > 0) {
			if matchPrefix, match := hasUrlOrDirPrefixUsers(h.authUrlsUsers, vhostReqPath, h.authDirsUsers, reqFsPath, userid); matchPrefix {
				success = match
			}
		}
		if success {
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
