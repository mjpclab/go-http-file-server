package serverHandler

import (
	"../i18n"
	"../util"
	"html/template"
	"net/http"
	"os"
	"path"
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
	DeleteUrl   string
}

type responseData struct {
	rawReqPath     string
	handlerReqPath string

	NeedAuth     bool
	AuthUserName string
	AuthSuccess  bool

	IsDownload bool
	IsUpload   bool
	IsMkdir    bool
	IsDelete   bool
	IsMutate   bool
	WantJson   bool

	CanUpload    bool
	CanMkdir     bool
	CanDelete    bool
	HasDeletable bool
	CanArchive   bool
	CanCors      bool

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
	AliasSubItems []os.FileInfo
	SubItemsHtml  []*itemHtml
	SubItemPrefix string
	SortState     SortState
	Context       *pathContext

	Lang  string
	Trans *i18n.Translation
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

func (h *handler) mergeAlias(
	rawRequestPath string,
	item os.FileInfo,
	subItems []os.FileInfo,
) (mergedSubItems, aliasSubItems []os.FileInfo, errs []error) {
	errs = []error{}

	if (item != nil && !item.IsDir()) || len(h.aliases) == 0 {
		return subItems, nil, errs
	}

	for _, alias := range h.aliases {
		subName, isChildAlias, ok := getAliasSubPart(alias, rawRequestPath)
		if !ok {
			continue
		}
		aliasCaseSensitive := alias.caseSensitive()

		var fsItem os.FileInfo
		if isChildAlias { // reached second-deepest path of alias
			var err error
			fsItem, err = os.Stat(alias.fsPath())
			if err != nil {
				errs = append(errs, err)
			}
		}

		matchExisted := false
		for i, subItem := range subItems {
			if !alias.namesEqual(subItem.Name(), subName) {
				continue
			}
			matchExisted = true
			if isVirtual(subItem) {
				continue
			}
			var baseItem os.FileInfo
			if fsItem != nil {
				baseItem = fsItem
			} else {
				baseItem = subItem
			}
			aliasSubItem := createVirtualFileInfo(subItem.Name(), baseItem, aliasCaseSensitive)
			aliasSubItems = append(aliasSubItems, aliasSubItem)
			subItems[i] = aliasSubItem
			if aliasCaseSensitive {
				break
			}
		}

		if !matchExisted {
			// fsItem could be nil
			aliasSubItem := createVirtualFileInfo(subName, fsItem, aliasCaseSensitive)
			aliasSubItems = append(aliasSubItems, aliasSubItem)
			subItems = append(subItems, aliasSubItem)
		}
	}

	return subItems, aliasSubItems, errs
}

func getSubItemPrefix(rawRequestPath string, tailSlash bool) string {
	if tailSlash {
		return "./"
	} else {
		return "./" + path.Base(rawRequestPath) + "/"
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
	if len(h.dirIndexes) == 0 {
		return
	}

	for _, index := range h.dirIndexes {
		for _, alias := range h.aliases {
			if !alias.isMatch(path.Clean(rawReqPath + "/" + index)) {
				continue
			}
			file, item, err = stat(alias.fsPath(), true)
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

	if baseItem == nil || !baseItem.IsDir() || h.emptyRoot {
		return
	}

	for _, index := range h.dirIndexes {
		file, item, err = stat(path.Clean(baseDir+"/"+index), true)
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

func (h *handler) getResponseData(r *http.Request) *responseData {
	var errs []error

	requestUri := r.URL.Path
	if len(requestUri) == 0 {
		requestUri = "/"
	}
	tailSlash := requestUri[len(requestUri)-1] == '/'

	rawReqPath := util.CleanUrlPath(requestUri)
	reqPath := util.CleanUrlPath(rawReqPath[len(h.urlPrefix):]) // strip url prefix path
	reqFsPath, _ := util.NormalizeFsPath(h.root + reqPath)

	needAuth := h.getNeedAuth(rawReqPath, reqFsPath)
	authUserName := ""
	authSuccess := true
	if needAuth {
		var _authErr error
		authUserName, authSuccess, _authErr = h.verifyAuth(r)
		if _authErr != nil {
			errs = append(errs, _authErr)
		}
	}

	rawQuery := r.URL.RawQuery
	isDownload := false
	isUpload := false
	isMkdir := false
	isDelete := false
	isMutate := false
	switch {
	case strings.HasPrefix(rawQuery, "download"):
		isDownload = true
	case strings.HasPrefix(rawQuery, "upload") && r.Method == http.MethodPost:
		isUpload = true
		isMutate = true
	case strings.HasPrefix(rawQuery, "mkdir"):
		isMkdir = true
		isMutate = true
	case strings.HasPrefix(r.URL.RawQuery, "delete"):
		isDelete = true
		isMutate = true
	}
	wantJson := strings.HasPrefix(rawQuery, "json") || strings.Contains(rawQuery, "&json")

	status := http.StatusOK
	isRoot := rawReqPath == "/"

	pathEntries := getPathEntries(rawReqPath, tailSlash)
	var rootRelPath string
	if len(pathEntries) > 0 {
		rootRelPath = pathEntries[0].Path
	} else {
		rootRelPath = "./"
	}

	file, item, _statErr := stat(reqFsPath, authSuccess && !h.emptyRoot)
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

	subItems, _readdirErr := readdir(file, item, authSuccess && !isMutate && needResponseBody(r.Method))
	if _readdirErr != nil {
		errs = append(errs, _readdirErr)
		status = http.StatusInternalServerError
	}

	subItems, aliasSubItems, _mergeErrs := h.mergeAlias(rawReqPath, item, subItems)
	if len(_mergeErrs) > 0 {
		errs = append(errs, _mergeErrs...)
		status = http.StatusInternalServerError
	}

	subItems = h.FilterItems(subItems)
	rawSortBy, sortState := sortInfos(subItems, rawQuery, h.defaultSort)

	if h.emptyRoot && status == http.StatusOK && r.RequestURI != "/" {
		status = http.StatusNotFound
	}

	subItemPrefix := getSubItemPrefix(rawReqPath, tailSlash)

	canUpload := h.getCanUpload(item, rawReqPath, reqFsPath)
	canMkdir := h.getCanMkdir(item, rawReqPath, reqFsPath)
	canDelete := h.getCanDelete(item, rawReqPath, reqFsPath)
	hasDeletable := canDelete && len(subItems) > len(aliasSubItems)
	canArchive := h.getCanArchive(subItems, rawReqPath, reqFsPath)
	canCors := h.getCanCors(rawReqPath, reqFsPath)

	context := &pathContext{
		download:    isDownload,
		sort:        rawSortBy,
		defaultSort: h.defaultSort,
	}

	return &responseData{
		rawReqPath:     rawReqPath,
		handlerReqPath: reqPath,

		NeedAuth:     needAuth,
		AuthUserName: authUserName,
		AuthSuccess:  authSuccess,

		IsDownload: isDownload,
		IsUpload:   isUpload,
		IsMkdir:    isMkdir,
		IsDelete:   isDelete,
		IsMutate:   isMutate,
		WantJson:   wantJson,

		CanUpload:    canUpload,
		CanMkdir:     canMkdir,
		CanDelete:    canDelete,
		HasDeletable: hasDeletable,
		CanArchive:   canArchive,
		CanCors:      canCors,

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
		AliasSubItems: aliasSubItems,
		SubItemsHtml:  nil,
		SubItemPrefix: subItemPrefix,
		SortState:     sortState,
		Context:       context,
	}
}
