package serverHandler

import (
	"net/http"
	"net/url"
)

func (h *aliasHandler) content(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Vary", h.vary)
	header.Set("X-Content-Type-Options", "nosniff")
	if data.IsDownload {
		header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(data.ItemName))
	}

	item := data.Item
	file := data.File

	http.ServeContent(w, r, item.Name(), item.ModTime(), file)
}
