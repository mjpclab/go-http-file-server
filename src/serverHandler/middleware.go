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

		NeedAuth:     data.NeedAuth,
		AuthUserName: data.AuthUserName,
		AuthSuccess:  data.AuthSuccess,

		RestrictAccess: data.RestrictAccess,
		AllowAccess:    data.AllowAccess,

		WantJson: data.wantJson,

		Status: data.Status,

		Item:     data.Item,
		SubItems: data.SubItems,

		Logger: h.logger,
	}

	for i := range h.postMiddlewares {
		result := h.postMiddlewares[i](w, r, context)
		if result == middleware.Outputted {
			return true
		} else if result == middleware.SkipRests {
			break
		}
	}

	return false
}
