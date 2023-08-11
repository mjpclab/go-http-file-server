package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"net/http"
)

func (h *aliasHandler) applyMiddlewares(mids []middleware.Middleware, w http.ResponseWriter, r *http.Request, data *responseData, fsPath string) (processed bool) {
	if len(mids) == 0 {
		return
	}

	context := &middleware.Context{
		PrefixReqPath: data.prefixReqPath,
		VhostReqPath:  data.rawReqPath,
		AliasReqPath:  data.handlerReqPath,
		AliasFsPath:   fsPath,
		AliasFsRoot:   h.root,

		WantJson: data.wantJson,

		RestrictAccess: data.RestrictAccess,
		AllowAccess:    data.AllowAccess,

		NeedAuth:     data.NeedAuth,
		AuthUserName: data.AuthUserName,
		AuthSuccess:  data.AuthSuccess,

		CanUpload:  &data.CanUpload,
		CanMkdir:   &data.CanMkdir,
		CanDelete:  &data.CanDelete,
		CanArchive: &data.CanArchive,

		Status: &data.Status,

		Users:  h.users,
		Logger: h.logger,
	}

	if data.File != nil {
		context.File = &data.File
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
