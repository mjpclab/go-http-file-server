package serverHandler

import (
	"../util"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type archiveCallback func(f *os.File, fInfo os.FileInfo, relPath string) error

func (h *handler) visitFs(
	initFsPath, rawRequestPath, relPath string,
	statFs bool,
	archiveCallback archiveCallback,
) {
	alias, hasAlias := h.aliases.byUrlPath(rawRequestPath)

	var fsPath string
	if hasAlias {
		fsPath = alias.fsPath
	} else {
		fsPath = initFsPath
	}

	if hasAlias && rawRequestPath != "/" {
		statFs = true
	}

	var fInfo os.FileInfo
	var childInfos []os.FileInfo
	if statFs {
		// wrap func to run defer ASAP
		err := func() error {
			f, err := os.Open(fsPath)
			if f != nil {
				defer f.Close()
			}
			h.errHandler.LogError(err)

			if err != nil {
				if os.IsExist(err) {
					return err
				}
				fInfo = newFakeFileInfo(path.Base(fsPath), true) // prefix path for alias
			} else {
				fInfo, err = f.Stat()
				if h.errHandler.LogError(err) {
					return err
				}
			}

			if len(relPath) > 0 {
				if err := archiveCallback(f, fInfo, relPath); err != nil {
					return err
				}
			}

			if f != nil && fInfo.IsDir() {
				childInfos, err = f.Readdir(0)
				if h.errHandler.LogError(err) {
					return err
				}
				childInfos = h.FilterItems(childInfos)
			}

			return nil
		}()
		if err != nil {
			return
		}
	} else {
		fInfo = newFakeFileInfo(path.Base(fsPath), true)
	}

	if fInfo.IsDir() {
		childAliases := map[string]string{}
		for _, alias := range h.aliases {
			if alias.isChildOf(rawRequestPath) {
				childAliases[alias.urlPath] = alias.fsPath
				continue
			}

			// just generate a path for walk down, mapped value is not important
			if alias.isSuccessorOf(rawRequestPath) {
				succPath := alias.urlPath[len(rawRequestPath):]
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

		// childInfo can be regular dir/file, or aliased item that shadows regular dir/file
		for _, childInfo := range childInfos {
			childPath := "/" + childInfo.Name()
			childFsPath := fsPath + childPath
			childRawRequestPath := util.CleanUrlPath(rawRequestPath + childPath)
			childRelPath := relPath + childPath

			if childAliasedFsPath, hasChildAlias := childAliases[childRawRequestPath]; hasChildAlias {
				h.visitFs(childAliasedFsPath, childRawRequestPath, childRelPath, statFs, archiveCallback)
				delete(childAliases, childRawRequestPath) // delete walked alias
			} else {
				h.visitFs(childFsPath, childRawRequestPath, childRelPath, statFs, archiveCallback)
			}
		}

		// rest aliases are standalone aliases that not shadows exists dir/file
		for childRawRequestPath, childAliasedFsPath := range childAliases {
			childRelPath := relPath + "/" + path.Base(childRawRequestPath)
			h.visitFs(childAliasedFsPath, childRawRequestPath, childRelPath, statFs, archiveCallback)
		}
	}
}

func (h *handler) archive(
	w http.ResponseWriter,
	r *http.Request,
	pageData *responseData,
	fileSuffix string,
	contentType string,
	cbWriteFile archiveCallback,
) {
	var itemName string
	_, hasAlias := h.aliases.byUrlPath(pageData.rawReqPath)
	if hasAlias {
		itemName = path.Base(pageData.rawReqPath)
	}
	if len(itemName) == 0 || itemName == "/" {
		itemName = pageData.ItemName
	}

	targetFilename := itemName + fileSuffix
	writeArchiveHeader(w, contentType, targetFilename)

	if !needResponseBody(r.Method) {
		return
	}

	h.visitFs(
		path.Clean(h.root+pageData.handlerReqPath),
		pageData.rawReqPath,
		"",
		!h.emptyRoot,
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			go h.logArchive(targetFilename, relPath, r)
			err := cbWriteFile(f, fInfo, relPath)
			h.errHandler.LogError(err)
			return err
		},
	)
}

func writeArchiveHeader(w http.ResponseWriter, contentType, filename string) {
	filename = url.PathEscape(filename)

	header := w.Header()
	header.Set("Content-Type", contentType)
	header.Set("Content-Disposition", "attachment; filename*=UTF-8''"+filename)
	header.Set("Cache-Control", "public, max-age=0")
	w.WriteHeader(http.StatusOK)
}
