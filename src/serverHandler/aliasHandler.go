package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"mjpclab.dev/ghfs/src/user"
	"mjpclab.dev/ghfs/src/util"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var defaultHandler = http.NotFoundHandler()

type aliasHandler struct {
	alias

	emptyRoot    bool
	autoDirSlash int
	hsts         bool
	hstsMaxAge   string
	toHttps      bool
	toHttpsPort  string // with prefix ":"
	defaultSort  string

	users  *user.List
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

	auth    *hierarchyAvailability
	index   *hierarchyAvailability
	upload  *hierarchyAvailability
	mkdir   *hierarchyAvailability
	delete  *hierarchyAvailability
	archive *hierarchyAvailability
	cors    *hierarchyAvailability

	globalRestrictAccess []string
	restrictAccessUrls   pathStringsList
	restrictAccessDirs   pathStringsList

	globalHeaders [][2]string
	headersUrls   pathHeadersList
	headersDirs   pathHeadersList

	postMiddlewares []middleware.Middleware
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
	session, data := h.getSessionData(r)
	h.logErrors(session.errors)
	if session.file != nil {
		defer func() {
			session.file.Close()
		}()
	}

	if !session.allowAccess {
		if !h.applyMiddlewares(h.postMiddlewares, w, r, session, data) {
			h.page(w, r, session, data)
		}
		return
	}

	if session.needAuth {
		h.notifyAuth(w)
	}

	if session.authSuccess {
		if session.requestAuth {
			returnUrl := h.extractNoAuthUrl(r, session, data)
			http.Redirect(w, r, returnUrl, http.StatusFound)
			return
		}

		if session.redirectAction == addSlashSuffix {
			redirectWithQuery(w, r, session.prefixReqPath+"/", h.autoDirSlash)
			return
		} else if session.redirectAction == removeSlashSuffix {
			redirectWithQuery(w, r, session.prefixReqPath[:len(session.prefixReqPath)-1], h.autoDirSlash)
			return
		}

		if data.CanCors {
			cors(w)
		}

		header(w, session.headers)

		if session.isMutate && h.mutate(w, r, session, data) {
			return
		} else if session.isArchive {
			switch session.archiveFormat {
			case tarFmt:
				if h.tar(w, r, session, data) {
					return
				}
			case tgzFmt:
				if h.tgz(w, r, session, data) {
					return
				}
			case zipFmt:
				if h.zip(w, r, session, data) {
					return
				}
			}
		}
	}

	if h.applyMiddlewares(h.postMiddlewares, w, r, session, data) {
		return
	}

	// final process
	if session.wantJson {
		h.json(w, r, session, data)
	} else if shouldServeAsContent(session.file, data.Item) {
		h.content(w, r, session, data)
	} else {
		h.page(w, r, session, data)
	}
}

func newAliasHandler(
	p *param.Param,
	vhostCtx *vhostContext,
	currentAlias alias,
	allAliases aliases,
) *aliasHandler {
	emptyRoot := p.EmptyRoot && currentAlias.url == "/"

	globalRestrictAccess := p.GlobalRestrictAccess
	globalRestrictAccess = vhostCtx.restrictAccessUrls.mergePrefixMatched(globalRestrictAccess, util.HasUrlPrefixDir, currentAlias.url)
	globalRestrictAccess = vhostCtx.restrictAccessDirs.mergePrefixMatched(globalRestrictAccess, util.HasFsPrefixDir, currentAlias.dir)
	globalRestrictAccess = util.InPlaceDedup(globalRestrictAccess)

	globalHeaders := p.GlobalHeaders
	globalHeaders = vhostCtx.headersUrls.mergePrefixMatched(globalHeaders, util.HasUrlPrefixDir, currentAlias.url)
	globalHeaders = vhostCtx.headersDirs.mergePrefixMatched(globalHeaders, util.HasFsPrefixDir, currentAlias.dir)

	h := &aliasHandler{
		alias:        currentAlias,
		emptyRoot:    emptyRoot,
		autoDirSlash: p.AutoDirSlash,
		hsts:         p.Hsts,
		hstsMaxAge:   strconv.Itoa(p.HstsMaxAge),
		toHttps:      p.ToHttps,
		toHttpsPort:  p.ToHttpsPort,
		defaultSort:  p.DefaultSort,

		users:  vhostCtx.users,
		theme:  vhostCtx.theme,
		logger: vhostCtx.logger,

		dirIndexes: p.DirIndexes,
		aliases:    allAliases.filterSuccessor(currentAlias.url),

		auth:    newHierarchyAvailability(currentAlias.url, currentAlias.dir, p.GlobalAuth, p.AuthUrls, vhostCtx.authUrlsUsers, p.AuthDirs, vhostCtx.authDirsUsers),
		index:   newHierarchyAvailability(currentAlias.url, currentAlias.dir, false, p.IndexUrls, vhostCtx.indexUrlsUsers, p.IndexDirs, vhostCtx.indexDirsUsers),
		upload:  newHierarchyAvailability(currentAlias.url, currentAlias.dir, p.GlobalUpload, p.UploadUrls, vhostCtx.uploadUrlsUsers, p.UploadDirs, vhostCtx.uploadDirsUsers),
		mkdir:   newHierarchyAvailability(currentAlias.url, currentAlias.dir, p.GlobalMkdir, p.MkdirUrls, vhostCtx.mkdirUrlsUsers, p.MkdirDirs, vhostCtx.mkdirDirsUsers),
		delete:  newHierarchyAvailability(currentAlias.url, currentAlias.dir, p.GlobalDelete, p.DeleteUrls, vhostCtx.deleteUrlsUsers, p.DeleteDirs, vhostCtx.deleteDirsUsers),
		archive: newHierarchyAvailability(currentAlias.url, currentAlias.dir, p.GlobalArchive, p.ArchiveUrls, vhostCtx.archiveUrlsUsers, p.ArchiveDirs, vhostCtx.archiveDirsUsers),
		cors:    newHierarchyAvailability(currentAlias.url, currentAlias.dir, p.GlobalCors, p.CorsUrls, nil, p.CorsDirs, nil),

		globalRestrictAccess: globalRestrictAccess,
		restrictAccessUrls:   vhostCtx.restrictAccessUrls.filterSuccessor(false, util.HasUrlPrefixDir, currentAlias.url),
		restrictAccessDirs:   vhostCtx.restrictAccessDirs.filterSuccessor(false, util.HasFsPrefixDir, currentAlias.dir),

		globalHeaders: globalHeaders,
		headersUrls:   vhostCtx.headersUrls.filterSuccessor(false, util.HasUrlPrefixDir, currentAlias.url),
		headersDirs:   vhostCtx.headersDirs.filterSuccessor(false, util.HasFsPrefixDir, currentAlias.dir),

		shows:     vhostCtx.shows,
		showDirs:  vhostCtx.showDirs,
		showFiles: vhostCtx.showFiles,
		hides:     vhostCtx.hides,
		hideDirs:  vhostCtx.hideDirs,
		hideFiles: vhostCtx.hideFiles,

		postMiddlewares: p.PostMiddlewares,
	}
	return h
}
