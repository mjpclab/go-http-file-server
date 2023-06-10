package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (h *aliasHandler) content(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Vary", h.vary)
	header.Set("X-Content-Type-Options", "nosniff")
	if data.IsDownload {
		header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(data.ItemName))
	}

	item := data.Item
	file := data.File

	if NeedResponseBody(r.Method) {
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
		return
	}

	header.Set("Accept-Ranges", "bytes")
	if lacksHeader(header, "Content-Type") {
		ctype, err := util.GetContentType(item.Name(), file)
		if err == nil {
			header.Set("Content-Type", ctype)
		} else if data.Status == http.StatusOK {
			data.Status = http.StatusInternalServerError
		}
	}
	header.Set("Content-Length", strconv.FormatInt(item.Size(), 10))
	header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	header.Set("Last-Modified", item.ModTime().UTC().Format(http.TimeFormat))

	w.WriteHeader(data.Status)
}
