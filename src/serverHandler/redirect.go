package serverHandler

import "net/http"

func (h *aliasHandler) redirectWithSlashSuffix(w http.ResponseWriter, r *http.Request, pathWithoutSlashSuffix string) {
	target := pathWithoutSlashSuffix + "/"
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, target, h.forceDirSlash)
}
