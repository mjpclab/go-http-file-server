package serverHandler

import "net/http"

func redirectWithQuery(w http.ResponseWriter, r *http.Request, path string, code int) {
	target := path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, code)
}
