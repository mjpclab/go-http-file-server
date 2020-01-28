package serverHandler

import "net/http"

func (h *handler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	if data.HasInternalError {
		w.WriteHeader(http.StatusInternalServerError)
	} else if data.HasNotFoundError {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if needResponseBody(r.Method) {
		updateSubsItemHtml(data.SubItems)
		err := h.template.Execute(w, data)
		h.errHandler.LogError(err)
	}
}
