package serverHandler

import (
	"../util"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) content(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	if data.IsDownload {
		header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+data.ItemName)
	}

	item := data.Item
	file := data.File

	if needResponseBody(r.Method) {
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
		return
	}

	ctype, err := util.GetContentType(item.Name(), file)
	if err == nil {
		header.Set("Content-Type", ctype)
		header.Set("Content-Length", strconv.FormatInt(item.Size(), 10))
		header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
		header.Set("Last-Modified", item.ModTime().UTC().Format(http.TimeFormat))
	} else if data.Status == http.StatusOK {
		data.Status = http.StatusInternalServerError
	}

	w.WriteHeader(data.Status)
}
