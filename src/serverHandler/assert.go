package serverHandler

import (
	"../tpl/assert"
	"net/http"
	"time"
)

var initTime = time.Now()

func (h *handler) assert(w http.ResponseWriter, r *http.Request, assertPath string) {
	content, ok := assert.Get(assertPath)
	if !ok {
		return
	}

	header := w.Header()
	header.Set("Content-Type", content.ContentType)
	if needResponseBody(r.Method) {
		http.ServeContent(w, r, assertPath, initTime, content.ReadSeeker)
	}
}
