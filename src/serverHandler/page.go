package serverHandler

import (
	"html/template"
	"mjpclab.dev/ghfs/src/acceptHeaders"
	"mjpclab.dev/ghfs/src/i18n"
	tplUtil "mjpclab.dev/ghfs/src/tpl/util"
	"mjpclab.dev/ghfs/src/util"
	"net/http"
)

const TypeDir = template.HTML("dir")
const TypeFile = template.HTML("file")

func updateSubItemsHtml(data *responseData) {
	length := len(data.SubItems)
	data.SubItemsHtml = make([]itemHtml, length)

	for i, info := range data.SubItems {
		name := info.Name()
		urlEscapedName := tplUtil.FormatFileUrl(name)

		var displayName template.HTML
		var typ template.HTML
		var url string
		var readableSize template.HTML

		if info.IsDir() {
			displayName = tplUtil.FormatFilename(name) + "/"
			typ = TypeDir
			url = data.SubItemPrefix + urlEscapedName + "/" + data.Context.QueryString()
		} else {
			displayName = tplUtil.FormatFilename(name)
			typ = TypeFile
			url = data.SubItemPrefix + urlEscapedName + data.Context.FileQueryString()
			readableSize = tplUtil.FormatSize(info.Size())
		}

		var deleteUrl string
		if data.CanDelete && !isVirtual(info) {
			deleteUrl = name
		}

		data.SubItemsHtml[i] = itemHtml{
			Type:        typ,
			Url:         url,
			DisplayName: displayName,
			DisplaySize: readableSize,
			DisplayTime: tplUtil.FormatTime(info.ModTime()),
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

func (h *aliasHandler) page(w http.ResponseWriter, r *http.Request, data *responseData) {
	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Type", "text/html; charset=utf-8")
	if lacksHeader(header, "Cache-Control") {
		header.Set("Cache-Control", "public, max-age=0")
	}
	if len(h.pageVary) > 0 {
		header.Set("Vary", h.pageVary)
	}

	updateTranslation(r, data)

	if !NeedResponseBody(r.Method) {
		w.WriteHeader(data.Status)
		return
	}

	w.WriteHeader(data.Status)

	updateSubItemsHtml(data)
	err := h.theme.RenderPage(w, data)
	h.logError(err)
}
