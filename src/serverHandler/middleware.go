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
		AliasFsRoot:   h.root,

		WantJson: session.wantJson,

		AllowAccess: session.allowAccess,

		NeedAuth:     session.needAuth,
		AuthUserName: session.authUserName,
		AuthSuccess:  session.authSuccess,

		CanUpload:  &data.CanUpload,
		CanMkdir:   &data.CanMkdir,
		CanDelete:  &data.CanDelete,
		CanArchive: &data.CanArchive,

		Status: &data.Status,

		Users:  h.users,
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
