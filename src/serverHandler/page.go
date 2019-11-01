package serverHandler

import "net/http"

func (h *handler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	if data.hasInternalError {
		w.WriteHeader(http.StatusInternalServerError)
	} else if data.hasNotFoundError {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err := h.template.Execute(w, data)
	h.errHandler.LogError(err)
}
