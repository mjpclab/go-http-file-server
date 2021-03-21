package serverHandler

import (
	"net/http"
	"time"
)

var initTime = time.Now()

func (h *handler) asset(w http.ResponseWriter, r *http.Request, assetPath string) {
	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Cache-Control", "public, max-age=3600")
	if needResponseBody(r.Method) {
		h.theme.RenderAsset(w, r, assetPath)
	}
}
