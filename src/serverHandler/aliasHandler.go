package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"mjpclab.dev/ghfs/src/user"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var defaultHandler = http.NotFoundHandler()

var createFileServer func(aliasUrl, aliasFs string) http.Handler

type pathStrings struct {
	path    string
	strings []string
}

type aliasParam struct {
	users  user.List
	theme  theme.Theme
	logger *serverLog.Logger

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp

	restrictAccess     bool
	restrictAccessUrls []pathStrings
	restrictAccessDirs []pathStrings

	headersUrls []pathHeaders
	headersDirs []pathHeaders

	vary string
}

type aliasHandler struct {
	root          string
	emptyRoot     bool
	forceDirSlash int
	hsts          bool
	hstsMaxAge    string
	toHttps       bool
	toHttpsPort   string // with prefix ":"
	defaultSort   string
	aliasPrefix   string

	users  user.List
	theme  theme.Theme
	logger *serverLog.Logger

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp

	dirIndexes []string
	aliases    aliases

	globalAuth bool
	authUrls   []string
	authDirs   []string

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

	vary string

	postMiddlewares []middleware.Middleware

	fileServer http.Handler
}

func (h *aliasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// hsts redirect
	if h.hsts && h.tryHsts(w, r) {
		return
	}

	// https redirect
	if h.toHttps && h.tryToHttps(w, r) {
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
		fileServer = createFileServer(currentAlias.url, currentAlias.fs)
	}

	h := &aliasHandler{
		root:          currentAlias.fs,
		emptyRoot:     emptyRoot,
		forceDirSlash: p.ForceDirSlash,
		hsts:          p.Hsts,
		hstsMaxAge:    strconv.Itoa(p.HstsMaxAge),
		toHttps:       p.ToHttps,
		toHttpsPort:   p.ToHttpsPort,
		defaultSort:   p.DefaultSort,
		aliasPrefix:   currentAlias.url,

		users:  ap.users,
		theme:  ap.theme,
		logger: ap.logger,

		dirIndexes: p.DirIndexes,
		aliases:    aliases,

		globalAuth: p.GlobalAuth,
		authUrls:   p.AuthUrls,
		authDirs:   p.AuthDirs,

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

		shows:     ap.shows,
		showDirs:  ap.showDirs,
		showFiles: ap.showFiles,
		hides:     ap.hides,
		hideDirs:  ap.hideDirs,
		hideFiles: ap.hideFiles,

		fileServer: fileServer,

		vary: ap.vary,

		postMiddlewares: p.PostMiddlewares,
	}
	return h
}
