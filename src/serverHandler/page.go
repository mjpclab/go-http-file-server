package serverHandler

import (
	tplutil "../tpl/util"
	"net/http"
)

func updateSubItemsHtml(data *responseData) {
	length := len(data.SubItems)

	data.SubItemsHtml = make([]*itemHtml, length)

	for i := 0; i < length; i++ {
		info := data.SubItems[i]
		name := info.Name()

		data.SubItemsHtml[i] = &itemHtml{
			IsDir:   info.IsDir(),
			Link:    name,
			Name:    tplutil.FormatFilename(name),
			Size:    tplutil.FormatSize(info.Size()),
			ModTime: tplutil.FormatTime(info.ModTime()),
		}
	}
}

func (h *handler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	w.WriteHeader(data.Status)

	if needResponseBody(r.Method) {
		updateSubItemsHtml(data)
		err := h.template.Execute(w, data)
		h.errHandler.LogError(err)
	}
}
