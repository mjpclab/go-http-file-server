package serverHandler

import "net/http"

func (h *aliasHandler) redirectWithSlashSuffix(w http.ResponseWriter, r *http.Request, pathWithoutSlashSuffix string) {
	target := pathWithoutSlashSuffix + "/"
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, target, h.forceDirSlash)
}

func (h *aliasHandler) redirectWithoutForceAuth(w http.ResponseWriter, r *http.Request, data *responseData) {
	returnUrl := r.Header.Get("Referer")
	if len(returnUrl) == 0 {
		returnUrl = data.prefixReqPath + data.Context.QueryString()
	}

	http.Redirect(w, r, returnUrl, http.StatusFound)
}
