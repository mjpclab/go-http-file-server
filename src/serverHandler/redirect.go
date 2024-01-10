package serverHandler

import "net/http"

func redirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	target := path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, code)
}

func (h *aliasHandler) redirectWithoutRequestAuth(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	returnUrl := r.Header.Get("Referer")
	if len(returnUrl) == 0 {
		returnUrl = session.prefixReqPath + data.Context.QueryString()
	}

	http.Redirect(w, r, returnUrl, http.StatusFound)
}
