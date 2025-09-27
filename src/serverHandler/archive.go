package serverHandler

import (
	"mjpclab.dev/ghfs/src/util"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type archiveCallback func(f *os.File, fInfo os.FileInfo, relPath string) error

func matchSelection(info os.FileInfo, selections []string) (match bool, childSelections []string) {
	if len(selections) == 0 {
		return true, nil
	}

	name := info.Name()
	for _, selName := range selections {
		if util.IsPathEqual(selName, name) {
			match = true
			continue
		}

		if !info.IsDir() {
			continue
		}
		slashIndex := strings.IndexByte(selName, '/')
		if slashIndex <= 0 {
			continue
		}

		selNamePart1 := selName[:slashIndex]
		if util.IsPathEqual(selNamePart1, name) {
			match = true
			childSel := selName[slashIndex+1:]
			if len(childSel) > 0 {
				childSelections = append(childSelections, childSel)
			}
			continue
		}
	}

	return
}

func (h *aliasHandler) visitTreeNode(
	r *http.Request,
	urlPath, fsPath, relPath string,
	statNode bool,
	childSelections []string,
	archiveCallback archiveCallback,
) {
	select {
	case <-r.Context().Done():
		return
	default:
		break
	}

	needAuth, _ := h.needAuth("", urlPath, fsPath)
	userId, _, err := h.verifyAuth(r, urlPath, fsPath)
	if needAuth && err != nil {
		return
	}

	var fInfo os.FileInfo
	var childInfos []os.FileInfo
	// wrap func to run defer ASAP
	err = func() error {
		var f *os.File
		var err error
		if statNode {
			f, err = os.Open(fsPath)
			if f != nil {
				defer f.Close()
			}

			if h.logError(err) {
				if os.IsExist(err) {
					return err
				}
				fInfo = createPlaceholderFileInfo(path.Base(fsPath), true) // prefix path for alias
			} else {
				fInfo, err = f.Stat()
				if h.logError(err) {
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
			if h.logError(err) {
				return err
			}
		}

		return nil
	}()
	if err != nil {
		return
	}

	if fInfo.IsDir() && h.index.match(urlPath, fsPath, userId) {
		childInfos, _, _ = h.mergeAlias(urlPath, fInfo, childInfos, true)
		childInfos = h.FilterItems(childInfos)

		// childInfo can be regular dir/file, or aliased item that shadows regular dir/file
		for _, childInfo := range childInfos {
			matchChild, childChildSelections := matchSelection(childInfo, childSelections)
			if !matchChild {
				continue
			}

			childPath := "/" + childInfo.Name()
			childFsPath := fsPath + childPath
			childRawReqPath := util.CleanUrlPath(urlPath + childPath)
			childRelPath := relPath + childPath

			if childAlias, hasChildAlias := h.aliases.byUrlPath(childRawReqPath); hasChildAlias {
				h.visitTreeNode(r, childRawReqPath, childAlias.dir, childRelPath, true, childChildSelections, archiveCallback)
			} else {
				h.visitTreeNode(r, childRawReqPath, childFsPath, childRelPath, statNode, childChildSelections, archiveCallback)
			}
		}
	}
}

func (h *aliasHandler) archiveFiles(
	w http.ResponseWriter,
	r *http.Request,
	session *sessionContext,
	data *responseData,
	selections []string,
	contentType string,
	cbWriteFile archiveCallback,
) {
	writeArchiveHeader(w, contentType, session.outFileName)

	if !NeedResponseBody(r.Method) {
		return
	}

	h.visitTreeNode(
		r,
		session.vhostReqPath,
		session.fsPath,
		"",
		data.Item != nil, // not empty root
		selections,
		func(f *os.File, fInfo os.FileInfo, relPath string) error {
			h.logArchive(session.outFileName, relPath, r)
			err := cbWriteFile(f, fInfo, relPath)
			h.logError(err)
			return err
		},
	)
}

func writeArchiveHeader(w http.ResponseWriter, contentType, filename string) {
	filename = url.PathEscape(filename)

	header := w.Header()
	header.Set("Content-Type", contentType)
	header.Set("Content-Disposition", "attachment; filename="+filename+"; filename*=UTF-8''"+filename)
	header.Set("Cache-Control", "public, max-age=0")
	w.WriteHeader(http.StatusOK)
}

func (h *aliasHandler) normalizeArchiveSelections(r *http.Request) ([]string, bool) {
	if h.logError(r.ParseForm()) {
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
