package serverHandler

import (
	"net/http"
)

func (h *handler) asset(w http.ResponseWriter, r *http.Request, assetPath string) {
	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Cache-Control", "public, max-age=3600")
	h.theme.RenderAsset(w, r, assetPath)
}
