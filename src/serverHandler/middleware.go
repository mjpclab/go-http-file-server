package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"net/http"
)

func (h *aliasHandler) applyMiddlewares(mids []middleware.Middleware, w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) (processed bool) {
	if len(mids) == 0 {
		return
	}

	context := &middleware.Context{
		PrefixReqPath: session.prefixReqPath,
		VhostReqPath:  session.vhostReqPath,
		AliasReqPath:  session.aliasReqPath,
		AliasFsPath:   session.fsPath,
		AliasFsRoot:   h.dir,

		WantJson: session.wantJson,

		AllowAccess: session.allowAccess,

		NeedAuth:     session.needAuth,
		AuthUserName: data.AuthUserName,
		AuthSuccess:  session.authSuccess,

		Status: &data.Status,

		Logger: h.logger,
	}

	if session.file != nil {
		context.File = &session.file
	}
	if data.Item != nil {
		context.FileInfo = &data.Item
	}

	for i := range mids {
		result := mids[i](w, r, context)
		if result == middleware.Outputted {
			return true
		} else if result == middleware.SkipRests {
			break
		}
	}

	return false
}
