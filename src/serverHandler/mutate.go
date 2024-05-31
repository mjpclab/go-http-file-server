package serverHandler

import (
	"mjpclab.dev/ghfs/src/shimgo"
	"net/http"
)

func (h *aliasHandler) mutate(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	if r.Method != shimgo.Net_Http_MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	success := false

	switch {
	case data.IsUpload:
		if data.CanUpload {
			success = h.saveUploadFiles(data.AuthUserName, h.fs+session.aliasReqPath, data.CanMkdir, data.CanDelete, data.AliasSubItems, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case data.IsMkdir:
		if data.CanMkdir && !h.logError(r.ParseForm()) {
			success = h.mkdirs(data.AuthUserName, h.fs+session.aliasReqPath, r.Form["name"], data.AliasSubItems, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case data.IsDelete:
		if data.CanDelete && !h.logError(r.ParseForm()) {
			success = h.deleteItems(data.AuthUserName, h.fs+session.aliasReqPath, r.Form["name"], data.AliasSubItems, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
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
		reqPath := session.prefixReqPath

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
