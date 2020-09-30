package serverHandler

import "net/http"

func (h *handler) header(w http.ResponseWriter) {
	if len(h.globalHeaders) == 0 {
		return
	}
	header := w.Header()
	for _, headerPair := range h.globalHeaders {
		header.Set(headerPair[0], headerPair[1])
	}
}
