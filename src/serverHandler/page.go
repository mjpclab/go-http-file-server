package serverHandler

import (
	"html/template"
	"mjpclab.dev/ghfs/src/acceptHeaders"
	"mjpclab.dev/ghfs/src/i18n"
	tplUtil "mjpclab.dev/ghfs/src/tpl/util"
	"mjpclab.dev/ghfs/src/util"
	"net/http"
)

const typeDir = template.HTML("dir")
const typeFile = template.HTML("file")

func updateSubItemsHtml(data *responseData) {
	length := len(data.SubItems)
	if length == 0 {
		return
	}
	data.SubItemsHtml = make([]itemHtml, length)

	dirSuffix := "/" + data.Context.QueryString()
	fileSuffix := data.Context.SubFileQueryString()

	for i, info := range data.SubItems {
		name := info.Name()

		var displayName template.HTML
		var typ template.HTML
		var url string
		var readableSize template.HTML

		if info.IsDir() {
			displayName = tplUtil.FormatFilename(name) + "/"
			typ = typeDir
			url = data.SubItemPrefix + tplUtil.FormatFileUrl(name) + dirSuffix
		} else {
			displayName = tplUtil.FormatFilename(name)
			typ = typeFile
			url = data.SubItemPrefix + tplUtil.FormatFileUrl(name) + fileSuffix
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

func (h *aliasHandler) page(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) {
	header := w.Header()
	header.Set("Vary", session.vary)
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Type", "text/html; charset=utf-8")
	if lacksHeader(header, "Cache-Control") {
		header.Set("Cache-Control", "public, max-age=0")
	}

	updateTranslation(r, data)

	w.WriteHeader(data.Status)

	if !NeedResponseBody(r.Method) {
		return
	}

	updateSubItemsHtml(data)
	err := h.theme.RenderPage(w, data)
	h.logError(err)
}
