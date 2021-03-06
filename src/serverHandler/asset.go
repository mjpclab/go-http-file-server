package serverHandler

import (
	"net/http"
	"time"
)

var initTime = time.Now()

func (h *handler) asset(w http.ResponseWriter, r *http.Request, assetPath string) {
	content, ok := h.theme.Assets.Get(assetPath)
	if !ok {
		return
	}

	header := w.Header()
	header.Set("Content-Type", content.ContentType)
	header.Set("Cache-Control", "public, max-age=3600")
	if needResponseBody(r.Method) {
		http.ServeContent(w, r, assetPath, initTime, content.ReadSeeker)
	}
}
