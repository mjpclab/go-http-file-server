package serverHandler

import (
	"net/http"
	"strings"
)

func (h *handler) mutate(w http.ResponseWriter, r *http.Request, data *responseData) {
	success := false

	switch {
	case data.IsUpload:
		if data.CanUpload && r.Method == http.MethodPost {
			success = h.saveUploadFiles(h.root+data.handlerReqPath, data.CanDelete, data.AliasSubItems, r)
		}
	case data.IsMkdir:
		if data.CanMkdir && !h.errHandler.LogError(r.ParseForm()) {
			success = h.mkdirs(h.root+data.handlerReqPath, r.Form["name"], data.AliasSubItems)
		}
	case data.IsDelete:
		if data.CanDelete && !h.errHandler.LogError(r.ParseForm()) {
			success = h.deleteItems(h.root+data.handlerReqPath, r.Form["name"], data.AliasSubItems)
		}
	}

	if data.WantJson {
		header := w.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")
		header.Set("Cache-Control", "public, max-age=0")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte{'{', '"', 's', 'u', 'c', 'c', 'e', 's', 's', '"', ':'})
		if success {
			w.Write([]byte{'t', 'r', 'u', 'e'})
		} else {
			w.Write([]byte{'f', 'a', 'l', 's', 'e'})
		}
		w.Write([]byte{'}'})
	} else {
		reqPath := r.RequestURI
		qsIndex := strings.IndexByte(reqPath, '?')
		if qsIndex >= 0 {
			reqPath = reqPath[:qsIndex]
		}
		http.Redirect(w, r, reqPath, http.StatusFound)
	}
}
