package serverHandler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type jsonItem struct {
	IsDir     bool      `json:"isDir"`
	IsVirtual bool      `json:"isVirtual"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"modTime"`
}

type jsonResponseData struct {
	NeedAuth           bool        `json:"needAuth"`
	AuthUserName       string      `json:"authUserName"`
	AuthSuccess        bool        `json:"authSuccess"`
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
		IsDir:     info.IsDir(),
		IsVirtual: isVirtual(info),
		Name:      info.Name(),
		Size:      info.Size(),
		ModTime:   info.ModTime(),
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
		NeedAuth:           data.NeedAuth,
		AuthUserName:       data.AuthUserName,
		AuthSuccess:        data.AuthSuccess,
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

func (h *aliasHandler) json(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	if !needResponseBody(r.Method) {
		w.WriteHeader(data.Status)
		return
	}

	var bodyW io.Writer
	if compressW, encoding, useCompressW := getCompressWriter(w, r); useCompressW {
		header.Set("Content-Encoding", encoding)
		bodyW = compressW
		defer compressW.Close()
	} else {
		bodyW = w
	}
	w.WriteHeader(data.Status)

	jsonData := getJsonData(data)
	encoder := json.NewEncoder(bodyW)
	err := encoder.Encode(jsonData)
	h.logError(err)
}
