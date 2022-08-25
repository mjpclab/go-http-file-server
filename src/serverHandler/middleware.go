package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"net/http"
)

func (h *aliasHandler) postMiddleware(w http.ResponseWriter, r *http.Request, data *responseData, fsPath string) (processed bool) {
	if len(h.postMiddlewares) == 0 {
		return false
	}

	context := &middleware.Context{
		PrefixReqPath: data.prefixReqPath,
		VhostReqPath:  data.rawReqPath,
		AliasReqPath:  data.handlerReqPath,
		AliasFsPath:   fsPath,
		AliasFsRoot:   h.root,

		Item:     data.Item,
		SubItems: data.SubItems,

		Status: data.Status,
	}

	for i := range h.postMiddlewares {
		if h.postMiddlewares[i](w, r, context) {
			return true
		}
	}

	return false
}
