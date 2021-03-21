package serverHandler

import (
	"../acceptHeaders"
	"../i18n"
	tplutil "../tpl/util"
	"../util"
	"html/template"
	"io"
	"net/http"
)

const TypeDir = template.HTML("dir")
const TypeFile = template.HTML("file")

func updateSubItemsHtml(data *responseData) {
	length := len(data.SubItems)
	data.SubItemsHtml = make([]*itemHtml, length)
	dirQueryString := data.Context.QueryString()

	for i, info := range data.SubItems {
		name := info.Name()
		urlEscapedName := tplutil.FormatFileUrl(name)

		var displayName template.HTML
		var typ template.HTML
		var url string
		var readableSize template.HTML

		if info.IsDir() {
			displayName = tplutil.FormatFilename(name) + "/"
			typ = TypeDir
			url = data.SubItemPrefix + urlEscapedName + "/" + dirQueryString
		} else {
			displayName = tplutil.FormatFilename(name)
			typ = TypeFile
			url = data.SubItemPrefix + urlEscapedName
			readableSize = tplutil.FormatSize(info.Size())
		}

		var deleteUrl string
		if data.CanDelete && !isVirtual(info) {
			deleteUrl = name
		}

		data.SubItemsHtml[i] = &itemHtml{
			Type:        typ,
			Url:         url,
			DisplayName: displayName,
			DisplaySize: readableSize,
			DisplayTime: tplutil.FormatTime(info.ModTime()),
			DeleteUrl:   deleteUrl,
		}
	}
}

func updateTranslation(r *http.Request, data *responseData) {
	accepts := acceptHeaders.ParseAccepts(util.AsciiToLowerCase(r.Header.Get("Accept-Language")))
	index, _, ok := accepts.GetPreferredValue(i18n.LanguageTags)
	if !ok {
		index = 0
	}
	data.Lang = i18n.LanguageTags[index]
	data.Trans = i18n.Dictionaries[index].Trans
}

func (h *handler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "public, max-age=0")

	updateTranslation(r, data)
	header.Set("Content-Language", data.Lang)

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

	updateSubItemsHtml(data)
	err := h.theme.RenderPage(bodyW, data)
	if err != nil {
		go h.errHandler.LogError(err)
	}
}
