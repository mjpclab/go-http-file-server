package serverHandler

import "net/http"

func (h *handler) content(w http.ResponseWriter, r *http.Request, data *responseData) {
	if needResponseBody(r.Method) {
		item := data.Item
		file := data.File
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
	}
}
