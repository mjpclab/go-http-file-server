package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"net/http"
)

func (h *aliasHandler) middleware(w http.ResponseWriter, r *http.Request, data *responseData, fsPath string) (processed bool) {
	if len(h.middlewares) == 0 {
		return false
	}

	context := &middleware.Context{
		PrefixReqPath: data.prefixReqPath,
		VhostReqPath:  data.rawReqPath,
		AliasReqPath:  data.handlerReqPath,
		AliasFsPath:   fsPath,
		AliasFsRoot:   h.root,
	}

	for i := range h.middlewares {
		if h.middlewares[i](w, r, context) {
			return true
		}
	}

	return false
}
