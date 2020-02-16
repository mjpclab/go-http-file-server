package serverHandler

import (
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func getContentType(item os.FileInfo, file *os.File) (string, error) {
	ext := path.Ext(item.Name())
	ctype := mime.TypeByExtension(ext)
	if len(ctype) > 0 {
		return ctype, nil
	}

	var buf [512]byte
	n, err := file.Read(buf[:])
	if err != nil {
		return ctype, err
	}

	ctype = http.DetectContentType(buf[:n])
	return ctype, nil
}

func (h *handler) content(w http.ResponseWriter, r *http.Request, data *responseData) {
	item := data.Item
	file := data.File

	if needResponseBody(r.Method) {
		http.ServeContent(w, r, item.Name(), item.ModTime(), file)
		return
	}

	ctype, err := getContentType(item, file)
	if err == nil {
		header := w.Header()
		header.Set("Content-Type", ctype)
		header.Set("Content-Length", strconv.FormatInt(item.Size(), 10))
		header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
		header.Set("Last-Modified", item.ModTime().UTC().Format(http.TimeFormat))
	} else {
		// take effect only if (!HasForbiddenError && !HasNotFoundError)
		data.HasInternalError = true
	}

	writeHeader(w, r, data)
}
