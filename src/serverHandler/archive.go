package serverHandler

import (
	"../util"
	"net/http"
	"net/url"
	"os"
	"path"
)

type archiveCallback func(f *os.File, fInfo os.FileInfo, relPath string) error

func (h *handler) visitFs(
	initFsPath, rawRequestPath, relPath string,
	statFs bool,
	archiveCallback archiveCallback,
) {
	var fsPath string
	alias, hasAlias := h.aliases.byUrlPath(rawRequestPath)
	if hasAlias {
		fsPath = alias.fsPath
		if alias.urlPath != "/" {
			statFs = true
		}
	} else {
		fsPath = initFsPath
	}

	var fInfo os.FileInfo
	var childInfos []os.FileInfo
	// wrap func to run defer ASAP
	err := func() error {
		var f *os.File
		var err error
		if statFs {
			f, err = os.Open(fsPath)
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
		} else {
			fInfo = newFakeFileInfo(path.Base(fsPath), true)
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
		}

		return nil
	}()
	if err != nil {
		return
	}

	if fInfo.IsDir() {
		childInfos, _, _ := h.mergeAlias(rawRequestPath, fInfo, childInfos)
		childInfos = h.FilterItems(childInfos)

		// childInfo can be regular dir/file, or aliased item that shadows regular dir/file
		for _, childInfo := range childInfos {
			childPath := "/" + childInfo.Name()
			childFsPath := fsPath + childPath
			childRawRequestPath := util.CleanUrlPath(rawRequestPath + childPath)
			childRelPath := relPath + childPath

			if childAlias, hasChildAlias := h.aliases.byUrlPath(childRawRequestPath); hasChildAlias {
				h.visitFs(childAlias.fsPath, childRawRequestPath, childRelPath, statFs, archiveCallback)
			} else {
				h.visitFs(childFsPath, childRawRequestPath, childRelPath, statFs, archiveCallback)
			}
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
