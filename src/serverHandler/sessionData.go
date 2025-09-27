package serverHandler

import (
	"html/template"
	"mjpclab.dev/ghfs/src/acceptHeaders"
	"mjpclab.dev/ghfs/src/i18n"
	"mjpclab.dev/ghfs/src/util"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type redirectAction int

const (
	noRedirect redirectAction = iota
	addSlashSuffix
	removeSlashSuffix
)

type archiveFormat int

const (
	tarFmt archiveFormat = iota
	tgzFmt
	zipFmt
)

const contentTypeJson = "application/json"

var acceptContentTypes = []string{
	"text/html",
	"application/xhtml+xml",
	"application/xml",
	contentTypeJson,
}
var acceptJsonIndex = len(acceptContentTypes) - 1

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

type sessionContext struct {
	prefixReqPath string
	vhostReqPath  string
	aliasReqPath  string
	fsPath        string

	allowAccess bool

	needAuth    bool
	requestAuth bool
	authUserId  int
	authSuccess bool

	redirectAction redirectAction
	vary           string
	headers        [][2]string

	wantJson bool

	isUpload bool
	isMkdir  bool
	isDelete bool
	isMutate bool

	isArchive     bool
	archiveFormat archiveFormat

	file *os.File

	query url.Values

	outFileName string

	errors []error
}

type responseData struct {
	AuthUserName string

	IsSimple   bool
	IsDownload bool

	CanIndex     bool
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
			fsItem, err = os.Stat(alias.dir)
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

func getItemName(info os.FileInfo, r *http.Request, prefixReqPath string) (itemName string) {
	itemName = path.Base(prefixReqPath)
	if len(itemName) > 0 && itemName != "/" {
		return
	}

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
			file, item, err = stat(alias.dir, true)
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
	fsPath := filepath.Clean(h.dir + reqPath)

	query := r.URL.Query()
	queryPrefix := getQueryPrefix(r.URL.RawQuery)

	status := http.StatusOK

	needAuth, requestAuth := h.needAuth(queryPrefix, vhostReqPath, fsPath)
	authUserId, authUserName, _authErr := h.verifyAuth(r, vhostReqPath, fsPath)
	authSuccess := !needAuth || _authErr == nil
	if !authSuccess {
		status = http.StatusUnauthorized
		errs = append(errs, _authErr)
	}

	headers := h.getHeaders(vhostReqPath, fsPath, authSuccess)

	accepts := acceptHeaders.ParseAccepts(r.Header.Get("Accept"))
	acceptIndex, _, _ := accepts.GetPreferredValue(acceptContentTypes)
	wantJson := acceptIndex == acceptJsonIndex

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

	canIndex := authSuccess && redirectAction == noRedirect && h.index.match(vhostReqPath, fsPath, authUserId)
	indexFile, indexItem, _statIdxErr := h.statIndexFile(vhostReqPath, fsPath, item, canIndex)
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
	vary := "accept, accept-encoding"
	if restrictAccess {
		vary += ", referer, origin"
	}
	if !allowAccess {
		status = http.StatusForbidden
	}

	canIndex = canIndex && allowAccess

	itemName := getItemName(item, r, prefixReqPath)

	var outFileName string

	isUpload := false
	isMkdir := false
	isDelete := false
	isMutate := false
	switch queryPrefix {
	case "upload":
		isUpload = true
		isMutate = true
	case "mkdir":
		isMkdir = true
		isMutate = true
	case "delete":
		isDelete = true
		isMutate = true
	}

	isArchive := false
	var arFmt archiveFormat
	switch queryPrefix {
	case "tar":
		isArchive = true
		arFmt = tarFmt
		outFileName = ".tar"
	case "tgz":
		isArchive = true
		arFmt = tgzFmt
		outFileName = ".tar.tz"
	case "zip":
		isArchive = true
		arFmt = zipFmt
		outFileName = ".zip"
	}
	if isArchive {
		arName, _ := getQueryValue(query, queryPrefix)
		if len(arName) > 0 {
			outFileName = arName
		} else {
			outFileName = itemName + outFileName
		}
	}

	_, isSimple := query["simple"]
	dlName, isDownload := getQueryValue(query, "download")
	if isDownload {
		if len(dlName) > 0 {
			outFileName = dlName
		} else {
			outFileName = itemName
		}
	}

	subItems, _readdirErr := readdir(file, item, canIndex && !isMutate && !isArchive && NeedResponseBody(r.Method))
	if _readdirErr != nil {
		errs = append(errs, _readdirErr)
		status = http.StatusInternalServerError
	}

	subItems, aliasSubItems, _mergeErrs := h.mergeAlias(vhostReqPath, item, subItems, canIndex)
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
	sort := query.Get("sort")
	if len(sort) > 2 {
		sort = sort[:2]
	}
	sortState := sortInfos(subItems, sort, h.defaultSort)

	if h.emptyRoot && status == http.StatusOK && len(vhostReqPath) > 1 {
		status = http.StatusNotFound
	}

	subItemPrefix := getSubItemPrefix(currDirRelPath, vhostReqPath, tailSlash)

	isDir := item != nil && item.IsDir()
	canUpload := allowAccess && authSuccess && isDir && h.upload.match(vhostReqPath, fsPath, authUserId)
	canMkdir := allowAccess && authSuccess && isDir && h.mkdir.match(vhostReqPath, fsPath, authUserId)
	canDelete := allowAccess && authSuccess && isDir && h.delete.match(vhostReqPath, fsPath, authUserId)
	hasDeletable := canDelete && len(subItems) > len(aliasSubItems)
	canArchive := allowAccess && authSuccess && h.archive.match(vhostReqPath, fsPath, authUserId)
	canCors := allowAccess && authSuccess && h.cors.match(vhostReqPath, fsPath, authUserId)
	loginAvail := len(authUserName) == 0 && h.users.Len() > 0

	session = &sessionContext{
		prefixReqPath: prefixReqPath,
		vhostReqPath:  vhostReqPath,
		aliasReqPath:  reqPath,
		fsPath:        fsPath,

		allowAccess: allowAccess,

		needAuth:    needAuth,
		requestAuth: requestAuth,
		authUserId:  authUserId,
		authSuccess: authSuccess,

		redirectAction: redirectAction,
		vary:           vary,
		headers:        headers,

		wantJson: wantJson,

		isUpload: isUpload,
		isMkdir:  isMkdir,
		isDelete: isDelete,
		isMutate: isMutate,

		isArchive:     isArchive,
		archiveFormat: arFmt,

		outFileName: outFileName,

		file: file,

		query: query,

		errors: errs,
	}
	data = &responseData{
		AuthUserName: authUserName,

		IsSimple:   isSimple,
		IsDownload: isDownload,

		CanIndex:     canIndex,
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
		SubItems:      subItems,
		AliasSubItems: aliasSubItems,
		SubItemsHtml:  nil,
		SubItemPrefix: subItemPrefix,
		SortState:     sortState,
		Context: pathContext{
			simple:      isSimple,
			download:    isDownload,
			sort:        sort,
			defaultSort: h.defaultSort,
		},
	}
	return
}
