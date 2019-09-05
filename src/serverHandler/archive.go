package serverHandler

import (
	"../serverError"
	"../util"
	"os"
	"path"
)

func (h *handler) visitFs(
	initFsPath, rawRequestPath, relPath string,
	callback func(*os.File, os.FileInfo, string),
) {
	aliasedFsPath, hasAlias := h.aliases[rawRequestPath]

	var fsPath string
	if hasAlias {
		fsPath = aliasedFsPath
	} else {
		fsPath = initFsPath
	}

	f, err := os.Open(fsPath)
	if serverError.LogError(err) {
		return
	}

	fInfo, err := f.Stat()
	if serverError.LogError(err) {
		return
	}

	if len(relPath) > 0 {
		callback(f, fInfo, relPath)
	}

	if fInfo.IsDir() {
		childAliases := map[string]string{}
		for urlPath, fsPath := range h.aliases {
			if path.Dir(urlPath) == rawRequestPath {
				childAliases[urlPath] = fsPath
			}
		}

		childInfos, err := f.Readdir(0)
		if serverError.LogError(err) {
			return
		}

		for _, childInfo := range childInfos {
			childPath := "/" + childInfo.Name()
			childFsPath := fsPath + childPath
			childRawRequestPath := util.CleanUrlPath(rawRequestPath + childPath)
			childRelPath := relPath + childPath

			if childAliasedFsPath, hasChildAlias := childAliases[childRawRequestPath]; hasChildAlias {
				h.visitFs(childAliasedFsPath, childRawRequestPath, childRelPath, callback)
				delete(childAliases, childRawRequestPath)
			} else {
				h.visitFs(childFsPath, childRawRequestPath, childRelPath, callback)
			}
		}

		for childRawRequestPath, childAliasedFsPath := range childAliases {
			childRelPath := relPath + "/" + path.Base(childRawRequestPath)
			h.visitFs(childAliasedFsPath, childRawRequestPath, childRelPath, callback)
		}
	}
}
