package serverHandler

import (
	"html/template"
	"mjpclab.dev/ghfs/src/i18n"
	"mjpclab.dev/ghfs/src/util"
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
	prefixReqPath  string
	rawReqPath     string
	handlerReqPath string

	NeedAuth     bool
	AuthUserName string
	AuthSuccess  bool

	RestrictAccess bool
	AllowAccess    bool

	Headers [][2]string

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

	IsRoot      bool
	Path        string
	Paths       []pathEntry
	RootRelPath string

	File          *os.File
	Item          os.FileInfo
	ItemName      string
	SubItems      []os.FileInfo
	AliasSubItems []os.FileInfo
	SubItemsHtml  []itemHtml
	SubItemPrefix string
	SortState     SortState
	Context       pathContext

	NeedDirSlashRedirect bool

	Lang  string
	Trans *i18n.Translation
}

func getPathEntries(currDirRelPath, path string, tailSlash bool) []pathEntry {
	pathSegs := make([]string, 1, 12)
	pathSegs[0] = "/"
	for refPath := path[1:]; len(refPath) > 0; {
		slashIndex := strings.IndexByte(refPath, '/')
		if slashIndex < 0 {
			pathSegs = append(pathSegs, refPath)
			break
		} else {
			pathSegs = append(pathSegs, refPath[:slashIndex])
			refPath = refPath[slashIndex+1:]
		}
	}

	pathCount := len(pathSegs)

	pathDepth := pathCount
	if !tailSlash {
		pathDepth--
	}

	pathEntries := make([]pathEntry, pathCount)
	for n := 1; n <= pathCount; n++ {
		var relPath string
		if n < pathDepth {
			relPath = strings.Repeat("../", pathDepth-n)
		} else if n == pathDepth {
			relPath = currDirRelPath
		} else /*if n == pathDepth+1*/ {
			relPath = currDirRelPath + pathSegs[pathDepth] + "/"
		}

		i := n - 1
		pathEntries[i] = pathEntry{
			Name: pathSegs[i],
			Path: relPath,
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

func (h *aliasHandler) mergeAlias(
	rawRequestPath string,
	item os.FileInfo,
	subItems []os.FileInfo,
	doMerge bool,
) (mergedSubItems, aliasSubItems []os.FileInfo, errs []error) {
	if !doMerge || (item != nil && !item.IsDir()) || len(h.aliases) == 0 {
		return subItems, nil, errs
	}

	for _, alias := range h.aliases {
		subName, noMore, ok := alias.nextPartOf(rawRequestPath)
		if !ok {
			continue
		}

		var fsItem os.FileInfo
		if noMore { // reached second-deepest path of alias
			var err error
			fsItem, err = os.Stat(alias.fs)
			if err != nil {
				errs = append(errs, err)
			}
		}

		matchExisted := false
		for i, subItem := range subItems {
			if !util.IsPathEqual(subItem.Name(), subName) {
				continue
			}
			matchExisted = true
			if isVirtual(subItem) {
				continue
			}
			aliasSubItem := createVirtualFileInfo(subItem.Name(), fsItem)
			aliasSubItems = append(aliasSubItems, aliasSubItem)
			subItems[i] = aliasSubItem
			break
		}

		if !matchExisted {
			// fsItem could be nil
			aliasSubItem := createVirtualFileInfo(subName, fsItem)
			aliasSubItems = append(aliasSubItems, aliasSubItem)
			subItems = append(subItems, aliasSubItem)
		}
	}

	return subItems, aliasSubItems, errs
}

func getCurrDirRelPath(reqPath, prefixReqPath string) string {
	if len(reqPath) == 1 && len(prefixReqPath) > 1 && prefixReqPath[len(prefixReqPath)-1] != '/' {
		return "./" + path.Base(prefixReqPath) + "/"
	} else {
		return "./"
	}
}

func getSubItemPrefix(currDirRelPath, rawRequestPath string, tailSlash bool) string {
	if tailSlash {
		return currDirRelPath
	} else {
		return currDirRelPath + path.Base(rawRequestPath) + "/"
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

func (h *aliasHandler) statIndexFile(rawReqPath, baseDir string, baseItem os.FileInfo, doStat bool) (file *os.File, item os.FileInfo, err error) {
	if !doStat || len(h.dirIndexes) == 0 {
		return
	}

	for _, index := range h.dirIndexes {
		for _, alias := range h.aliases {
			if !alias.isMatch(path.Clean(rawReqPath + "/" + index)) {
				continue
			}
			file, item, err = stat(alias.fs, true)
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

func dereferenceSymbolLinks(reqFsPath string, subItems []os.FileInfo) (errs []error) {
	baseFsPath := reqFsPath + "/"

	for i := range subItems {
		if subItems[i].Mode()&os.ModeSymlink != 0 {
			dereferencedItem, err := os.Stat(baseFsPath + subItems[i].Name())
			if err != nil {
				errs = append(errs, err)
			} else {
				subItems[i] = dereferencedItem
			}
		}
	}

	return
}

func (h *aliasHandler) getResponseData(r *http.Request) (data *responseData, fsPath string) {
	var errs []error

	prefixReqPath := r.URL.RawPath // init by pathTransformHandler
	rawReqPath := r.URL.Path
	tailSlash := rawReqPath[len(rawReqPath)-1] == '/'

	reqPath := util.CleanUrlPath(rawReqPath[len(h.aliasPrefix):])
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

	headers := h.getHeaders(rawReqPath, reqFsPath, authSuccess)

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

	currDirRelPath := getCurrDirRelPath(rawReqPath, prefixReqPath)
	pathEntries := getPathEntries(currDirRelPath, rawReqPath, tailSlash)
	var rootRelPath string
	if len(pathEntries) > 0 {
		rootRelPath = pathEntries[0].Path
	} else {
		rootRelPath = currDirRelPath
	}

	file, item, _statErr := stat(reqFsPath, authSuccess && !h.emptyRoot)
	if _statErr != nil {
		errs = append(errs, _statErr)
		status = getStatusByErr(_statErr)
	}

	needDirSlashRedirect := h.forceDirSlash > 0 && prefixReqPath[len(prefixReqPath)-1] != '/' && item != nil && item.IsDir()

	indexFile, indexItem, _statIdxErr := h.statIndexFile(rawReqPath, reqFsPath, item, authSuccess && !needDirSlashRedirect)
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

	allowAccess := h.isAllowAccess(r, rawReqPath, reqFsPath, file, item)

	itemName := getItemName(item, r)

	subItems, _readdirErr := readdir(file, item, authSuccess && !isMutate && !needDirSlashRedirect && allowAccess && NeedResponseBody(r.Method))
	if _readdirErr != nil {
		errs = append(errs, _readdirErr)
		status = http.StatusInternalServerError
	}

	subItems, aliasSubItems, _mergeErrs := h.mergeAlias(rawReqPath, item, subItems, authSuccess && !needDirSlashRedirect && allowAccess)
	if len(_mergeErrs) > 0 {
		errs = append(errs, _mergeErrs...)
		status = http.StatusInternalServerError
	}

	_dereferenceErrs := dereferenceSymbolLinks(reqFsPath, subItems)
	if len(_dereferenceErrs) > 0 {
		errs = append(errs, _dereferenceErrs...)
	}

	// update `needDirSlashRedirect` for dangling intermediate alias directory
	if !needDirSlashRedirect && h.forceDirSlash > 0 && len(subItems) > 0 && prefixReqPath[len(prefixReqPath)-1] != '/' {
		needDirSlashRedirect = true
	}

	subItems = h.FilterItems(subItems)
	rawSortBy, sortState := sortInfos(subItems, rawQuery, h.defaultSort)

	if h.emptyRoot && status == http.StatusOK && len(rawReqPath) > 1 {
		status = http.StatusNotFound
	}

	subItemPrefix := getSubItemPrefix(currDirRelPath, rawReqPath, tailSlash)

	canUpload := h.getCanUpload(item, rawReqPath, reqFsPath)
	canMkdir := h.getCanMkdir(item, rawReqPath, reqFsPath)
	canDelete := h.getCanDelete(item, rawReqPath, reqFsPath)
	hasDeletable := canDelete && len(subItems) > len(aliasSubItems)
	canArchive := h.getCanArchive(subItems, rawReqPath, reqFsPath)
	canCors := h.getCanCors(rawReqPath, reqFsPath)

	context := pathContext{
		download:    isDownload,
		sort:        rawSortBy,
		defaultSort: h.defaultSort,
	}

	return &responseData{
		prefixReqPath:  prefixReqPath,
		rawReqPath:     rawReqPath,
		handlerReqPath: reqPath,

		NeedAuth:     needAuth,
		AuthUserName: authUserName,
		AuthSuccess:  authSuccess,

		RestrictAccess: h.restrictAccess,
		AllowAccess:    allowAccess,

		Headers: headers,

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

		IsRoot:      isRoot,
		Path:        rawReqPath,
		Paths:       pathEntries,
		RootRelPath: rootRelPath,

		File:          file,
		Item:          item,
		ItemName:      itemName,
		SubItems:      subItems,
		AliasSubItems: aliasSubItems,
		SubItemsHtml:  nil,
		SubItemPrefix: subItemPrefix,
		SortState:     sortState,
		Context:       context,

		NeedDirSlashRedirect: needDirSlashRedirect,
	}, reqFsPath
}
