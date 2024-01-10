package serverHandler

import (
	"net/http"
	"strings"
)

func (h *aliasHandler) mutate(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	success := false

	switch {
	case data.IsUpload:
		if data.CanUpload {
			success = h.saveUploadFiles(session.authUserName, h.root+session.aliasReqPath, data.CanMkdir, data.CanDelete, data.AliasSubItems, r)
		}
	case data.IsMkdir:
		if data.CanMkdir && !h.logError(r.ParseForm()) {
			success = h.mkdirs(session.authUserName, h.root+session.aliasReqPath, r.Form["name"], data.AliasSubItems, r)
		}
	case data.IsDelete:
		if data.CanDelete && !h.logError(r.ParseForm()) {
			success = h.deleteItems(session.authUserName, h.root+session.aliasReqPath, r.Form["name"], data.AliasSubItems, r)
		}
	}

	if session.wantJson {
		header := w.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")
		header.Set("Cache-Control", "public, max-age=0")

		if success {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success":false}`))
		}
	} else {
		reqPath := r.RequestURI
		qsIndex := strings.IndexByte(reqPath, '?')
		if qsIndex >= 0 {
			reqPath = reqPath[:qsIndex]
		}

		ctxQsList := r.Form["contextquerystring"]
		ctxQsListLen := len(ctxQsList)
		if ctxQsListLen > 0 {
			ctxQs := ctxQsList[ctxQsListLen-1]
			if len(ctxQs) > 0 {
				reqPath += ctxQs
			}
		}

		if success {
			http.Redirect(w, r, reqPath, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
