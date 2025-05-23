package serverHandler

import (
	"net/http"
)

func (h *aliasHandler) mutate(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) (ok bool) {
	if r.Method != http.MethodPost {
		data.Status = http.StatusMethodNotAllowed
		return
	}

	switch {
	case session.isUpload:
		if data.CanUpload {
			ok = h.saveUploadFiles(data.AuthUserName, h.dir+session.aliasReqPath, data.CanMkdir, data.CanDelete, data.AliasSubItems, r)
		} else {
			data.Status = http.StatusBadRequest
			return
		}
	case session.isMkdir:
		if data.CanMkdir && !h.logError(r.ParseForm()) {
			ok = h.mkdirs(data.AuthUserName, h.dir+session.aliasReqPath, r.Form["name"], data.AliasSubItems, r)
		} else {
			data.Status = http.StatusBadRequest
			return
		}
	case session.isDelete:
		if data.CanDelete && !h.logError(r.ParseForm()) {
			ok = h.deleteItems(data.AuthUserName, h.dir+session.aliasReqPath, r.Form["name"], data.AliasSubItems, r)
		} else {
			data.Status = http.StatusBadRequest
			return
		}
	}

	if session.wantJson {
		header := w.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")

		if ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success":false}`))
		}
		return true
	}

	if ok {
		reqPath := session.prefixReqPath
		ctxQsList := r.Form["contextquerystring"]
		ctxQsListLen := len(ctxQsList)
		if ctxQsListLen > 0 {
			ctxQs := ctxQsList[ctxQsListLen-1]
			if len(ctxQs) > 0 {
				reqPath += ctxQs
			}
		}
		http.Redirect(w, r, reqPath, http.StatusFound)
		return
	}

	data.Status = http.StatusInternalServerError
	return
}
