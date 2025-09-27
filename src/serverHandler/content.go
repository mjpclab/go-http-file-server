package serverHandler

import (
	"net/http"
	"net/url"
)

func (h *aliasHandler) content(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	header := w.Header()
	header.Set("Vary", session.vary)
	header.Set("X-Content-Type-Options", "nosniff")
	if data.IsDownload {
		filename := url.PathEscape(session.outFileName)
		header.Set("Content-Disposition", "attachment; filename="+filename+"; filename*=UTF-8''"+filename)
	}

	item := data.Item
	file := session.file

	http.ServeContent(w, r, item.Name(), item.ModTime(), file)
}
