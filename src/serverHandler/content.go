package serverHandler

import (
	"mjpclab.dev/ghfs/src/shimgo"
	"mjpclab.dev/ghfs/src/util"
	"net/http"
	"os"
	"strconv"
	"time"
)

var serveContent = func(h *aliasHandler, w http.ResponseWriter, r *http.Request, info os.FileInfo, file *os.File) {
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

func (h *aliasHandler) content(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	if r.ProtoMajor <= 1 {
		if len(h.contentVaryV1) > 0 {
			header.Set("Vary", h.contentVaryV1)
		}
	} else {
		if len(h.contentVary) > 0 {
			header.Set("Vary", h.contentVary)
		}
	}
	if data.IsDownload {
		header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+shimgo.Net_Url_PathEscape(data.ItemName))
	}

	item := data.Item
	file := data.File

	if NeedResponseBody(r.Method) {
		serveContent(h, w, r, item, file)
		return
	}

	ctype, err := util.GetContentType(item.Name(), file)
	if err == nil {
		header.Set("Accept-Ranges", "bytes")
		if len(header.Get("Content-Type")) == 0 {
			header.Set("Content-Type", ctype)
		}
		header.Set("Content-Length", strconv.FormatInt(item.Size(), 10))
		header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
		header.Set("Last-Modified", item.ModTime().UTC().Format(http.TimeFormat))
	} else if data.Status == http.StatusOK {
		data.Status = http.StatusInternalServerError
	}

	w.WriteHeader(data.Status)
}
