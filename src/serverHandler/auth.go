package serverHandler

import (
	"errors"
	"net/http"
	"net/url"
)

const authQueryParam = "auth"

func (h *aliasHandler) needAuth(queryPrefix, vhostReqPath, reqFsPath string) (needAuth, requestAuth bool) {
	if queryPrefix == authQueryParam {
		return true, true
	}

	if h.auth.global {
		return true, false
	}

	if hasUrlOrDirPrefix(h.auth.urls, vhostReqPath, h.auth.dirs, reqFsPath) {
		return true, false
	}

	if matchPath, _ := hasUrlOrDirPrefixUsers(h.auth.urlsUsers, vhostReqPath, h.auth.dirsUsers, reqFsPath, -1); matchPath {
		return true, false
	}

	return false, false
}

func (h *aliasHandler) notifyAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"files\"")
}

func (h *aliasHandler) verifyAuth(r *http.Request, vhostReqPath, reqFsPath string) (authUserId int, authUserName string, err error) {
	inputUser, inputPass, hasAuthReq := r.BasicAuth()

	if hasAuthReq {
		userid, username, success := h.users.Auth(inputUser, inputPass)
		if success && userid >= 0 && (len(h.auth.urlsUsers) > 0 || len(h.auth.dirsUsers) > 0) {
			if matchPrefix, match := hasUrlOrDirPrefixUsers(h.auth.urlsUsers, vhostReqPath, h.auth.dirsUsers, reqFsPath, userid); matchPrefix {
				success = match
			}
		}
		if success {
			return userid, username, nil
		}
		err = errors.New(r.RemoteAddr + " auth failed")
	} else {
		err = errors.New(r.RemoteAddr + " missing auth info")
	}

	authUserId = -1
	return
}

func (h *aliasHandler) extractNoAuthUrl(r *http.Request, session *sessionContext, data *responseData) string {
	returnUrl, hasReturnUrl := getQueryValue(session.query, authQueryParam)
	if hasReturnUrl && len(returnUrl) > 0 {
		var err error
		returnUrl, err = url.QueryUnescape(returnUrl)
		if err == nil {
			return returnUrl
		}
	}

	returnUrl = r.Header.Get("Referer")
	if len(returnUrl) > 0 {
		return returnUrl
	}

	return session.prefixReqPath + data.Context.QueryString()
}
