package serverHandler

import "net/http"

func (h *handler) header(w http.ResponseWriter) {
	header := w.Header()

	for i := range h.globalHeaders {
		header.Add(h.globalHeaders[i][0], h.globalHeaders[i][1])
	}
}
