package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl"
	"mjpclab.dev/ghfs/src/user"
	"net/http"
	"regexp"
	"strings"
)

var defaultHandler = http.NotFoundHandler()

var createFileServer func(root string) http.Handler

type pathStrings struct {
	path    string
	strings []string
}

type aliasParam struct {
	users  user.List
	theme  tpl.Theme
	logger *serverLog.Logger

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp

	headersUrls []pathHeaders
	headersDirs []pathHeaders

	restrictAccess     bool
	restrictAccessUrls []pathStrings
	restrictAccessDirs []pathStrings

	pageVaryV1    string
	pageVary      string
	contentVaryV1 string
	contentVary   string
}

type aliasHandler struct {
	root          string
	emptyRoot     bool
	forceDirSlash int
	globalHsts    bool
	globalHttps   bool
	httpsPort     string // with prefix ":"
	defaultSort   string
	aliasPrefix   string

	users  user.List
	theme  tpl.Theme
	logger *serverLog.Logger

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp

	dirIndexes []string
	aliases    aliases

	restrictAccess       bool
	globalRestrictAccess []string
	restrictAccessUrls   []pathStrings
	restrictAccessDirs   []pathStrings

	globalHeaders [][2]string
	headersUrls   []pathHeaders
	headersDirs   []pathHeaders

	globalUpload bool
	uploadUrls   []string
	uploadDirs   []string

	globalMkdir bool
	mkdirUrls   []string
	mkdirDirs   []string

	globalDelete bool
	deleteUrls   []string
	deleteDirs   []string

	globalArchive bool
	archiveUrls   []string
	archiveDirs   []string

	globalCors bool
	corsUrls   []string
	corsDirs   []string

	globalAuth bool
	authUrls   []string
	authDirs   []string

	pageVaryV1    string
	pageVary      string
	contentVaryV1 string
	contentVary   string

	postMiddlewares []middleware.Middleware

	fileServer http.Handler
}

func (h *aliasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// hsts redirect
	if h.globalHsts && h.hsts(w, r) {
		return
	}

	// https redirect
	if h.globalHttps && h.https(w, r) {
		return
	}

	// asset
	const assetPrefix = "asset="
	if strings.HasPrefix(r.URL.RawQuery, assetPrefix) {
		assetPath := r.URL.RawQuery[len(assetPrefix):]
		h.asset(w, r, assetPath)
		return
	}

	// data
	data, fsPath := h.getResponseData(r)
	h.logErrors(data.errors)
	file := data.File
	if file != nil {
		defer file.Close()
	}

	if data.NeedAuth {
		h.needAuth(w, r)
	}
	if !data.AuthSuccess {
		if !h.postMiddleware(w, r, data, fsPath) {
			h.authFailed(w, data.Status)
		}
		return
	}

	if !data.AllowAccess {
		if !h.postMiddleware(w, r, data, fsPath) {
			h.accessRestricted(w, data.Status)
		}
		return
	}

	if data.NeedDirSlashRedirect {
		h.redirectWithSlashSuffix(w, r, data.prefixReqPath)
		return
	}

	header(w, data.Headers)

	if data.CanCors {
		cors(w)
	}

	if data.IsMutate {
		h.mutate(w, r, data)
		return
	}

	// archive
	if len(r.URL.RawQuery) >= 3 {
		switch r.URL.RawQuery[:3] {
		case "tar":
			h.tar(w, r, data)
			return
		case "tgz":
			h.tgz(w, r, data)
			return
		case "zip":
			h.zip(w, r, data)
			return
		}
	}

	if h.postMiddleware(w, r, data, fsPath) {
		return
	}

	// final process
	item := data.Item
	if data.WantJson {
		h.json(w, r, data)
	} else if shouldServeAsContent(file, item) {
		h.content(w, r, data)
	} else {
		h.page(w, r, data)
	}
}

func newAliasHandler(
	p *param.Param,
	ap *aliasParam,
	currentAlias alias,
	allAliases aliases,
) http.Handler {
	emptyRoot := p.EmptyRoot && currentAlias.url == "/"

	aliases := aliases{}
	for _, alias := range allAliases {
		if alias.isSuccessorOf(currentAlias.url) {
			aliases = append(aliases, alias)
		}
	}

	var fileServer http.Handler
	if !emptyRoot && createFileServer != nil { // for WSL 1 fix
		fileServer = createFileServer(currentAlias.fs)
	}

	h := &aliasHandler{
		root:          currentAlias.fs,
		emptyRoot:     emptyRoot,
		forceDirSlash: p.ForceDirSlash,
		globalHsts:    p.GlobalHsts,
		globalHttps:   p.GlobalHttps,
		httpsPort:     p.HttpsPort,
		defaultSort:   p.DefaultSort,
		aliasPrefix:   currentAlias.url,

		users:  ap.users,
		theme:  ap.theme,
		logger: ap.logger,

		dirIndexes: p.DirIndexes,
		aliases:    aliases,

		restrictAccess:       ap.restrictAccess,
		globalRestrictAccess: p.GlobalRestrictAccess,
		restrictAccessUrls:   ap.restrictAccessUrls,
		restrictAccessDirs:   ap.restrictAccessDirs,

		globalHeaders: p.GlobalHeaders,
		headersUrls:   ap.headersUrls,
		headersDirs:   ap.headersDirs,

		globalUpload: p.GlobalUpload,
		uploadUrls:   p.UploadUrls,
		uploadDirs:   p.UploadDirs,

		globalMkdir: p.GlobalMkdir,
		mkdirUrls:   p.MkdirUrls,
		mkdirDirs:   p.MkdirDirs,

		globalDelete: p.GlobalDelete,
		deleteUrls:   p.DeleteUrls,
		deleteDirs:   p.DeleteDirs,

		globalArchive: p.GlobalArchive,
		archiveUrls:   p.ArchiveUrls,
		archiveDirs:   p.ArchiveDirs,

		globalCors: p.GlobalCors,
		corsUrls:   p.CorsUrls,
		corsDirs:   p.CorsDirs,

		globalAuth: p.GlobalAuth,
		authUrls:   p.AuthUrls,
		authDirs:   p.AuthDirs,

		shows:     ap.shows,
		showDirs:  ap.showDirs,
		showFiles: ap.showFiles,
		hides:     ap.hides,
		hideDirs:  ap.hideDirs,
		hideFiles: ap.hideFiles,

		fileServer: fileServer,

		pageVaryV1:    ap.pageVaryV1,
		pageVary:      ap.pageVary,
		contentVaryV1: ap.contentVaryV1,
		contentVary:   ap.contentVary,

		postMiddlewares: p.PostMiddlewares,
	}
	return h
}
