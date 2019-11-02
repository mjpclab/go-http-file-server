package serverHandler

import (
	"../util"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func (h *handler) visitFs(
	initFsPath, rawRequestPath, relPath string,
	callback func(*os.File, os.FileInfo, string) error,
) {
	aliasedFsPath, hasAlias := h.aliases[rawRequestPath]

	var fsPath string
	if hasAlias {
		fsPath = aliasedFsPath
	} else {
		fsPath = initFsPath
	}

	f, err := os.Open(fsPath)
	h.errHandler.LogError(err)

	var fInfo os.FileInfo
	if err != nil {
		if os.IsExist(err) {
			return
		}
		fInfo = newFakeFileInfo(path.Base(fsPath), true)
	} else {
		fInfo, err = f.Stat()
		if h.errHandler.LogError(err) {
			return
		}
	}

	if len(relPath) > 0 {
		if callback(f, fInfo, relPath) != nil {
			return
		}
	}

	if fInfo.IsDir() {
		childAliases := map[string]string{}
		for aliasUrlPath, aliasFsPath := range h.aliases {
			if path.Dir(aliasUrlPath) == rawRequestPath {
				childAliases[aliasUrlPath] = aliasFsPath
				continue
			}

			if len(aliasUrlPath) > len(rawRequestPath) && util.HasUrlPrefixDir(aliasUrlPath, rawRequestPath) {
				succPath := aliasUrlPath[len(rawRequestPath):]
				if succPath[0] == '/' {
					succPath = succPath[1:]
				}
				childName := succPath[:strings.Index(succPath, "/")]
				childUrlPath := util.CleanUrlPath(rawRequestPath + "/" + childName)
				childFsPath := fsPath + "/" + childName
				childAliases[childUrlPath] = childFsPath
				continue
			}
		}

		var childInfos []os.FileInfo
		if f != nil {
			childInfos, err = f.Readdir(0)
			if h.errHandler.LogError(err) {
				return
			}
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

func writeArchiveHeader(w http.ResponseWriter, contentType, filename string) {
	filename = url.PathEscape(filename)

	header := w.Header()
	header.Set("Content-Type", contentType)
	header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+filename)
	header.Set("Cache-Control", "public, max-age=0")
	w.WriteHeader(http.StatusOK)
}
