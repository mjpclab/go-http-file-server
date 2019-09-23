package serverHandler

import (
	"../tpl/assert"
	"net/http"
	"time"
)

var initTime = time.Now()

func (h *handler) assert(w http.ResponseWriter, r *http.Request, assertPath string) {
	content := assert.Get(assertPath)

	header := w.Header()
	header.Set("Content-Type", content.ContentType)
	http.ServeContent(w, r, assertPath, initTime, content.ReadSeeker)
}
