package serverHandler

import (
	"../util"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type pathEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type itemHtml struct {
	Type        template.HTML
	Url         string
	DisplayName template.HTML
	DisplaySize template.HTML
	DisplayTime template.HTML
}

type responseData struct {
	rawReqPath     string
	handlerReqPath string

	errors []error
	Status int

	IsRoot        bool
	Path          string
	Paths         []*pathEntry
	RootRelPath   string
	File          *os.File
	Item          os.FileInfo
	ItemName      string
	SubItems      []os.FileInfo
	SubItemsHtml  []*itemHtml
	SubItemPrefix string

	CanUpload  bool
	CanArchive bool
	CanCors    bool
	NeedAuth   bool
}

func isSlash(c rune) bool {
	return c == '/'
}

func getPathEntries(path string, tailSlash bool) []*pathEntry {
	paths := []string{"/"}
	paths = append(paths, strings.FieldsFunc(path, isSlash)...)

	displayPathsCount := len(paths)

	pathsCount := displayPathsCount
	if !tailSlash {
		pathsCount--
	}

	pathEntries := make([]*pathEntry, displayPathsCount)
	for i := 0; i < displayPathsCount; i++ {
		var rPath string
		switch {
		case i < pathsCount-1:
			rPath = strings.Repeat("../", pathsCount-1-i)
		case i == pathsCount-1:
			rPath = "./"
		default:
			rPath = "./" + strings.Join(paths[pathsCount:], "/") + "/"
		}

		pathEntries[i] = &pathEntry{
			Name: paths[i],
			Path: rPath,
		}
	}

	return pathEntries
}

func stat(reqFsPath string, visitFs bool) (file *os.File, item os.FileInfo, err error) {
	if !visitFs {
		return
	}

	file, err = os.Open(reqFsPath)
	if err != nil {
		return
	}

	item, err = file.Stat()
	if err != nil {
		return
	}

	return
}

func readdir(file *os.File, item os.FileInfo, visitFs bool) (subItems []os.FileInfo, err error) {
	if !visitFs || file == nil || item == nil || !item.IsDir() {
		return
	}

	return file.Readdir(0)
}

func (h *handler) mergeAlias(rawRequestPath string, item os.FileInfo, subItems []os.FileInfo) ([]os.FileInfo, []error) {
	errs := []error{}

	if (item != nil && !item.IsDir()) || len(h.aliases) == 0 {
		return subItems, errs
	}

	for _, alias := range h.aliases {
		aliasUrlPath := alias.urlPath
		aliasFsPath := alias.fsPath

		if len(aliasUrlPath) <= len(rawRequestPath) {
			continue
		}

		suffixIndex := len(rawRequestPath)
		aliasPrefix := aliasUrlPath[:suffixIndex]
		aliasSuffix := aliasUrlPath[suffixIndex:]

		if aliasPrefix != rawRequestPath {
			continue
		}

		if len(aliasPrefix) != 1 && aliasSuffix[0] != '/' {
			// aliasUrlPath doesn't cover the whole directory name
			// e.g:
			// rawReqPath == "/abc/def/ghi"
			// aliasPrefix == "/abc/de"
			continue
		}
		if aliasSuffix[0] == '/' {
			aliasSuffix = aliasSuffix[1:]
		}

		slashIndex := strings.Index(aliasSuffix, "/")
		var nextName string
		if slashIndex >= 0 {
			nextName = aliasSuffix[:slashIndex]
		} else {
			nextName = aliasSuffix
		}

		var aliasSubItem os.FileInfo
		if path.Dir(aliasUrlPath) == rawRequestPath { // reached second deepest path of alias
			var err error
			aliasSubItem, err = os.Stat(aliasFsPath)
			if err == nil {
				aliasSubItem = newRenamedFileInfo(nextName, aliasSubItem)
			} else {
				errs = append(errs, err)
				aliasSubItem = newFakeFileInfo(nextName, true)
			}
		} else {
			aliasSubItem = newFakeFileInfo(nextName, true)
		}

		replaced := false
		for i, subItem := range subItems {
			if subItem.Name() == nextName {
				subItems[i] = aliasSubItem
				replaced = true
				break
			}
		}

		if !replaced {
			subItems = append(subItems, aliasSubItem)
		}
	}

	return subItems, errs
}

func getSubItemPrefix(requestPath string, tailSlash bool) string {
	if tailSlash {
		return "./"
	} else {
		return "./" + path.Base(requestPath) + "/"
	}
}

func getItemName(info os.FileInfo, r *http.Request) (itemName string) {
	if info != nil {
		itemName = info.Name()
	}
	if len(itemName) == 0 || itemName == "." || itemName == "/" {
		itemName = strings.Replace(r.Host, ":", "_", -1)
	}
	return
}

func sortSubItems(subInfos []os.FileInfo) {
	names := make([][]byte, len(subInfos))
	for i := range subInfos {
		names[i] = []byte(subInfos[i].Name())
	}

	sort.Slice(
		subInfos,
		func(prevIndex, nextIndex int) bool {
			prevInfo := subInfos[prevIndex]
			nextInfo := subInfos[nextIndex]

			prevIsDir := prevInfo.IsDir()
			nextIsDir := nextInfo.IsDir()

			if prevIsDir != nextIsDir {
				return prevIsDir
			}

			return util.CompareNumInStr(names[prevIndex], names[nextIndex])
		},
	)
}

func getStatusByErr(err error) int {
	switch {
	case os.IsPermission(err):
		return http.StatusForbidden
	case os.IsNotExist(err):
		return http.StatusNotFound
	case err != nil:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

func (h *handler) stateIndexFile(rawReqPath, baseDir string, baseItem os.FileInfo) (file *os.File, item os.FileInfo, err error) {
	if baseItem == nil || !baseItem.IsDir() || h.emptyRoot || len(h.dirIndexes) == 0 {
		return
	}

	for _, index := range h.dirIndexes {
		for _, alias := range h.aliases {
			if rawReqPath+"/"+index != alias.urlPath {
				continue
			}
			file, item, err = stat(alias.fsPath, true)
			if err != nil && file != nil {
				file.Close()
			}
			if err != nil && os.IsNotExist(err) {
				continue
			} else {
				return
			}
		}
	}

	for _, index := range h.dirIndexes {
		file, item, err = stat(baseDir+"/"+index, true)
		if err != nil && file != nil {
			file.Close()
		}
		if err != nil && os.IsNotExist(err) {
			continue
		} else {
			return
		}
	}

	return nil, nil, nil
}

func (h *handler) getResponseData(r *http.Request) (data *responseData) {
	requestUri := r.URL.Path
	tailSlash := requestUri[len(requestUri)-1] == '/'

	rawReqPath := util.CleanUrlPath(requestUri)
	reqPath := util.CleanUrlPath(rawReqPath[len(h.urlPrefix):]) // strip url prefix path
	errs := []error{}
	status := http.StatusOK
	isRoot := rawReqPath == "/"

	pathEntries := getPathEntries(rawReqPath, tailSlash)
	var rootRelPath string
	if len(pathEntries) > 0 {
		rootRelPath = pathEntries[0].Path
	} else {
		rootRelPath = "./"
	}

	reqFsPath, _absErr := util.NormalizeFsPath(h.root + reqPath)
	if _absErr != nil {
		reqFsPath = filepath.Clean(h.root + reqPath)
	}

	file, item, _statErr := stat(reqFsPath, !h.emptyRoot)
	if _statErr != nil {
		errs = append(errs, _statErr)
		status = getStatusByErr(_statErr)
	}

	indexFile, indexItem, _statIdxErr := h.stateIndexFile(rawReqPath, reqFsPath, item)
	if _statIdxErr != nil {
		errs = append(errs, _statIdxErr)
		status = getStatusByErr(_statIdxErr)
	} else if indexFile != nil {
		if indexItem != nil {
			file.Close()
			file = indexFile
			item = indexItem
		} else {
			indexFile.Close()
		}
	}

	itemName := getItemName(item, r)

	subItems, _readdirErr := readdir(file, item, needResponseBody(r.Method))
	if _readdirErr != nil {
		errs = append(errs, _readdirErr)
		status = http.StatusInternalServerError
	}

	subItems, _mergeErrs := h.mergeAlias(rawReqPath, item, subItems)
	if len(_mergeErrs) > 0 {
		errs = append(errs, _mergeErrs...)
		status = http.StatusInternalServerError
	}

	if h.emptyRoot && status == http.StatusOK && r.RequestURI != "/" {
		status = http.StatusNotFound
	}

	subItems = h.FilterItems(subItems)
	sortSubItems(subItems)

	subItemPrefix := getSubItemPrefix(reqPath, tailSlash)

	canUpload := h.getCanUpload(item, rawReqPath, reqFsPath)
	canArchive := h.getCanArchive(subItems, rawReqPath, reqFsPath)
	canCors := h.getCanCors(rawReqPath, reqFsPath)
	needAuth := h.getNeedAuth(rawReqPath, reqFsPath)

	data = &responseData{
		rawReqPath:     rawReqPath,
		handlerReqPath: reqPath,

		errors: errs,
		Status: status,

		IsRoot:        isRoot,
		Path:          rawReqPath,
		Paths:         pathEntries,
		RootRelPath:   rootRelPath,
		File:          file,
		Item:          item,
		ItemName:      itemName,
		SubItems:      subItems,
		SubItemsHtml:  nil,
		SubItemPrefix: subItemPrefix,

		CanUpload:  canUpload,
		CanArchive: canArchive,
		CanCors:    canCors,
		NeedAuth:   needAuth,
	}

	return
}
