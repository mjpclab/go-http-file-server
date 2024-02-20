package serverHandler

import (
	"html/template"
	"mjpclab.dev/ghfs/src/i18n"
	"mjpclab.dev/ghfs/src/util"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

type redirectAction int

const (
	noRedirect redirectAction = iota
	addSlashSuffix
	removeSlashSuffix
)

type sessionContext struct {
	prefixReqPath string
	vhostReqPath  string
	aliasReqPath  string
	fsPath        string

	allowAccess bool

	needAuth     bool
	requestAuth  bool
	authUserId   int
	authUserName string
	authSuccess  bool

	redirectAction redirectAction
	vary           string
	headers        [][2]string

	wantJson bool

	file *os.File

	errors []error
}

type responseData struct {
	IsDownload     bool
	IsDownloadFile bool
	IsUpload       bool
	IsMkdir        bool
	IsDelete       bool
	IsMutate       bool

	CanUpload    bool
	CanMkdir     bool
	CanDelete    bool
	HasDeletable bool
	CanArchive   bool
	CanCors      bool
	LoginAvail   bool

	Status int

	IsRoot      bool
	Path        string
	Paths       []pathEntry
	RootRelPath string

	Item          os.FileInfo
	ItemName      string
	SubItems      []os.FileInfo
	AliasSubItems []os.FileInfo
	SubItemsHtml  []itemHtml
	SubItemPrefix string
	SortState     SortState
	Context       pathContext

	Lang  string
	Trans *i18n.Translation
}

func getPathEntries(currDirRelPath, path string, tailSlash bool) (pathEntries []pathEntry, rootRelPath string) {
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

	pathEntries = make([]pathEntry, pathCount)
	for i := 0; i < pathCount; i++ {
		depth := i + 1
		var relPath string
		if depth < pathDepth {
			if i == 0 {
				relPath = strings.Repeat("../", pathDepth-depth)
			} else {
				// optimization: use existing string instead of generating new one
				// should got same result as above `if` block
				relPath = pathEntries[i-1].Path[3:]
			}
		} else if depth == pathDepth {
			relPath = currDirRelPath
		} else /*if depth == pathDepth+1*/ {
			relPath = currDirRelPath + pathSegs[pathDepth] + "/"
			if !tailSlash {
				relPath = relPath[:len(relPath)-1]
			}
		}

		pathEntries[i] = pathEntry{
			Name: pathSegs[i],
			Path: relPath,
		}
	}
	rootRelPath = pathEntries[0].Path

	return
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

func (h *aliasHandler) getSessionData(r *http.Request) (session *sessionContext, data *responseData) {
	var errs []error

	prefixReqPath := r.URL.RawPath // init by pathTransformHandler
	vhostReqPath := r.URL.Path
	tailSlash := vhostReqPath[len(vhostReqPath)-1] == '/'

	reqPath := util.CleanUrlPath(vhostReqPath[len(h.url):])
	fsPath := filepath.Clean(h.fs + reqPath)

	rawQuery := r.URL.RawQuery

	status := http.StatusOK

	needAuth, requestAuth := h.needAuth(rawQuery, vhostReqPath, fsPath)
	authUserId, authUserName, _authErr := h.verifyAuth(r, needAuth, vhostReqPath, fsPath)
	authSuccess := _authErr == nil
	if needAuth && !authSuccess {
		errs = append(errs, _authErr)
	}
	if !authSuccess {
		status = http.StatusUnauthorized
	}

	headers := h.getHeaders(vhostReqPath, fsPath, authSuccess)

	isDownload := false
	isDownloadFile := false
	isUpload := false
	isMkdir := false
	isDelete := false
	isMutate := false
	switch {
	case strings.HasPrefix(rawQuery, "downloadfile"):
		isDownload = true
		isDownloadFile = true
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

	isRoot := vhostReqPath == "/"

	currDirRelPath := getCurrDirRelPath(vhostReqPath, prefixReqPath)
	pathEntries, rootRelPath := getPathEntries(currDirRelPath, vhostReqPath, tailSlash)

	file, item, _statErr := stat(fsPath, authSuccess && !h.emptyRoot)
	if _statErr != nil {
		errs = append(errs, _statErr)
		status = getStatusByErr(_statErr)
	}

	redirectAction := noRedirect
	if h.autoDirSlash > 0 && len(vhostReqPath) > 1 && item != nil {
		if item.IsDir() {
			if prefixReqPath[len(prefixReqPath)-1] != '/' {
				redirectAction = addSlashSuffix
			}
		} else {
			if prefixReqPath[len(prefixReqPath)-1] == '/' {
				redirectAction = removeSlashSuffix
			}
		}
	}

	indexFile, indexItem, _statIdxErr := h.statIndexFile(vhostReqPath, fsPath, item, authSuccess && redirectAction == noRedirect)
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

	restrictAccess, allowAccess := h.isAllowAccess(r, vhostReqPath, fsPath, file, item)
	vary := h.vary
	if restrictAccess {
		vary += ", referer, origin"
	}
	if !allowAccess {
		status = http.StatusForbidden
	}

	itemName := getItemName(item, r)

	subItems, _readdirErr := readdir(file, item, allowAccess && authSuccess && !isMutate && redirectAction == noRedirect && NeedResponseBody(r.Method))
	if _readdirErr != nil {
		errs = append(errs, _readdirErr)
		status = http.StatusInternalServerError
	}

	subItems, aliasSubItems, _mergeErrs := h.mergeAlias(vhostReqPath, item, subItems, allowAccess && authSuccess && redirectAction == noRedirect)
	if len(_mergeErrs) > 0 {
		errs = append(errs, _mergeErrs...)
		status = http.StatusInternalServerError
	}

	_dereferenceErrs := dereferenceSymbolLinks(fsPath, subItems)
	if len(_dereferenceErrs) > 0 {
		errs = append(errs, _dereferenceErrs...)
	}

	// set `redirectAction` to `addSlashSuffix` for dangling intermediate alias directory
	if redirectAction == noRedirect && h.autoDirSlash > 0 && len(subItems) > 0 && prefixReqPath[len(prefixReqPath)-1] != '/' {
		redirectAction = addSlashSuffix
	}

	subItems = h.FilterItems(subItems)
	rawSortBy, sortState := sortInfos(subItems, rawQuery, h.defaultSort)

	if h.emptyRoot && status == http.StatusOK && len(vhostReqPath) > 1 {
		status = http.StatusNotFound
	}

	subItemPrefix := getSubItemPrefix(currDirRelPath, vhostReqPath, tailSlash)

	canUpload := allowAccess && authSuccess && h.getCanUpload(item, vhostReqPath, fsPath)
	canMkdir := allowAccess && authSuccess && h.getCanMkdir(item, vhostReqPath, fsPath)
	canDelete := allowAccess && authSuccess && h.getCanDelete(item, vhostReqPath, fsPath)
	hasDeletable := canDelete && len(subItems) > len(aliasSubItems)
	canArchive := allowAccess && authSuccess && h.getCanArchive(subItems, vhostReqPath, fsPath)
	canCors := allowAccess && authSuccess && h.getCanCors(vhostReqPath, fsPath)
	loginAvail := len(authUserName) == 0 && h.users.Len() > 0

	context := pathContext{
		download:     isDownload,
		downloadfile: isDownloadFile,
		sort:         rawSortBy,
		defaultSort:  h.defaultSort,
	}

	session = &sessionContext{
		prefixReqPath: prefixReqPath,
		vhostReqPath:  vhostReqPath,
		aliasReqPath:  reqPath,
		fsPath:        fsPath,

		allowAccess: allowAccess,

		needAuth:     needAuth,
		requestAuth:  requestAuth,
		authUserId:   authUserId,
		authUserName: authUserName,
		authSuccess:  authSuccess,

		redirectAction: redirectAction,
		vary:           vary,
		headers:        headers,

		wantJson: wantJson,

		file: file,

		errors: errs,
	}
	data = &responseData{
		IsDownload:     isDownload,
		IsDownloadFile: isDownloadFile,
		IsUpload:       isUpload,
		IsMkdir:        isMkdir,
		IsDelete:       isDelete,
		IsMutate:       isMutate,

		CanUpload:    canUpload,
		CanMkdir:     canMkdir,
		CanDelete:    canDelete,
		HasDeletable: hasDeletable,
		CanArchive:   canArchive,
		CanCors:      canCors,
		LoginAvail:   loginAvail,

		Status: status,

		IsRoot:      isRoot,
		Path:        vhostReqPath,
		Paths:       pathEntries,
		RootRelPath: rootRelPath,

		Item:          item,
		ItemName:      itemName,
		SubItems:      subItems,
		AliasSubItems: aliasSubItems,
		SubItemsHtml:  nil,
		SubItemPrefix: subItemPrefix,
		SortState:     sortState,
		Context:       context,
	}
	return
}
