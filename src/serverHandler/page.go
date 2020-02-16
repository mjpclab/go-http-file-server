package serverHandler

import "net/http"

func (h *handler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	writeHeader(w, r, data)

	if needResponseBody(r.Method) {
		updateSubsItemHtml(data.SubItems)
		err := h.template.Execute(w, data)
		h.errHandler.LogError(err)
	}
}
