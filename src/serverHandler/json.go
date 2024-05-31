package serverHandler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type jsonItem struct {
	IsDir   bool      `json:"isDir"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type jsonResponseData struct {
	IsRoot             bool        `json:"isRoot"`
	Path               string      `json:"path"`
	Paths              []pathEntry `json:"paths"`
	SubItemPrefix      string      `json:"subItemPrefix"`
	ContextQueryString string      `json:"contextQueryString"`
	CanUpload          bool        `json:"canUpload"`
	CanMkdir           bool        `json:"canMkdir"`
	CanDelete          bool        `json:"canDelete"`
	CanArchive         bool        `json:"canArchive"`
	CanCors            bool        `json:"canCors"`

	Item     *jsonItem   `json:"item"`
	SubItems []*jsonItem `json:"subItems"`
}

func getJsonItem(info os.FileInfo) *jsonItem {
	return &jsonItem{
		IsDir:   info.IsDir(),
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}
}

func getJsonData(data *responseData) *jsonResponseData {
	var item *jsonItem
	var subItems []*jsonItem

	if data.Item != nil {
		item = getJsonItem(data.Item)
	}

	subItems = make([]*jsonItem, len(data.SubItems))
	for i := range data.SubItems {
		subItems[i] = getJsonItem(data.SubItems[i])
	}

	return &jsonResponseData{
		IsRoot:             data.IsRoot,
		Path:               data.Path,
		Paths:              data.Paths,
		SubItemPrefix:      data.SubItemPrefix,
		ContextQueryString: data.Context.QueryString(),
		CanUpload:          data.CanUpload,
		CanMkdir:           data.CanMkdir,
		CanDelete:          data.CanDelete,
		CanArchive:         data.CanArchive,
		CanCors:            data.CanCors,

		Item:     item,
		SubItems: subItems,
	}
}

func (h *aliasHandler) json(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	header := w.Header()
	header.Set("Vary", session.vary)
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Type", "application/json; charset=utf-8")
	if lacksHeader(header, "Cache-Control") {
		header.Set("Cache-Control", "public, max-age=0")
	}

	if !NeedResponseBody(r.Method) {
		w.WriteHeader(data.Status)
		return
	}

	w.WriteHeader(data.Status)

	jsonData := getJsonData(data)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(jsonData)
	h.logError(err)
}
