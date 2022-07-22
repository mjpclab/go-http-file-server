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

func matchSelection(info os.FileInfo, selections []string) (matchName, matchPrefix bool, childSelections []string) {
	if len(selections) == 0 {
		return true, false, nil
	}

	name := info.Name()
	for _, selName := range selections {
		if util.IsPathEqual(selName, name) {
			matchName = true
			continue
		}

		slashIndex := strings.IndexByte(selName, '/')
		if slashIndex <= 0 {
			continue
		}

		selNamePart1 := selName[:slashIndex]
		if util.IsPathEqual(selNamePart1, name) {
			childSel := selName[slashIndex+1:]
			if len(childSel) > 0 {
				matchPrefix = true
				childSelections = append(childSelections, childSel)
			}
			continue
		}
	}

	return
}

func (h *handler) visitTreeNode(
	fsPath, rawReqPath, relPath string,
	statNode bool,
	childSelections []string,
	archiveCallback archiveCallback,
) {
	var fInfo os.FileInfo
	var childInfos []os.FileInfo
	// wrap func to run defer ASAP
	err := func() error {
		var f *os.File
		var err error
		if statNode {
			f, err = os.Open(fsPath)
			if f != nil {
				defer f.Close()
			}

			if h.errHandler.LogError(err) {
				if os.IsExist(err) {
					return err
				}
				fInfo = createPlaceholderFileInfo(path.Base(fsPath), true) // prefix path for alias
			} else {
				fInfo, err = f.Stat()
				if h.errHandler.LogError(err) {
					return err
				}
			}
		} else {
			fInfo = createPlaceholderFileInfo(path.Base(fsPath), true)
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
		childInfos, _, _ := h.mergeAlias(rawReqPath, fInfo, childInfos, true)
		childInfos = h.FilterItems(childInfos)

		// childInfo can be regular dir/file, or aliased item that shadows regular dir/file
		for _, childInfo := range childInfos {
			matchChildName, matchChildPrefix, childChildSelections := matchSelection(childInfo, childSelections)
			if !matchChildName && !matchChildPrefix {
				continue
			}

			childPath := "/" + childInfo.Name()
			childFsPath := fsPath + childPath
			childRawReqPath := util.CleanUrlPath(rawReqPath + childPath)
			childRelPath := relPath + childPath

			if childAlias, hasChildAlias := h.aliases.byUrlPath(childRawReqPath); hasChildAlias {
				h.visitTreeNode(childAlias.fs, childRawReqPath, childRelPath, true, childChildSelections, archiveCallback)
			} else {
				h.visitTreeNode(childFsPath, childRawReqPath, childRelPath, statNode, childChildSelections, archiveCallback)
			}
		}
	}
}

func (h *handler) archive(
	w http.ResponseWriter,
	r *http.Request,
	pageData *responseData,
	selections []string,
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

	h.visitTreeNode(
		path.Clean(h.root+pageData.handlerReqPath),
		pageData.rawReqPath,
		"",
		pageData.Item != nil, // not empty root
		selections,
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			h.logArchive(targetFilename, relPath, r)
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

func (h *handler) normalizeArchiveSelections(r *http.Request) ([]string, bool) {
	if h.errHandler.LogError(r.ParseForm()) {
		return nil, false
	}
	inputs := r.Form["name"]
	if len(inputs) == 0 {
		return nil, true
	}

	count := len(inputs)
	selections := make([]string, count)
	for i := 0; i < count; i++ {
		var ok bool
		selections[i], ok = getCleanDirFilePath(inputs[i])
		if !ok {
			h.logger.LogErrorString("archive: illegal path " + inputs[i])
			return nil, false
		}
	}

	return selections, true
}
