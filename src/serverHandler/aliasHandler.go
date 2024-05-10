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

	globalAuth    bool
	authUrls      []string
	authUrlsUsers pathIntsList
	authDirs      []string
	authDirsUsers pathIntsList

	globalUpload    bool
	uploadUrls      []string
	uploadUrlsUsers pathIntsList
	uploadDirs      []string
	uploadDirsUsers pathIntsList

	globalMkdir    bool
	mkdirUrls      []string
	mkdirUrlsUsers pathIntsList
	mkdirDirs      []string
	mkdirDirsUsers pathIntsList

	globalDelete    bool
	deleteUrls      []string
	deleteUrlsUsers pathIntsList
	deleteDirs      []string
	deleteDirsUsers pathIntsList

	globalArchive bool
	archiveUrls   []string
	archiveDirs   []string

	globalCors bool
	corsUrls   []string
	corsDirs   []string

	globalRestrictAccess []string
	restrictAccessUrls   pathStringsList
	restrictAccessDirs   pathStringsList

	globalHeaders [][2]string
	headersUrls   pathHeadersList
	headersDirs   pathHeadersList

	vary string

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
			h.redirectWithoutRequestAuth(w, r, session, data)
			return
		}

		if session.redirectAction == addSlashSuffix {
			redirect(w, r, session.prefixReqPath+"/", h.autoDirSlash)
			return
		} else if session.redirectAction == removeSlashSuffix {
			redirect(w, r, session.prefixReqPath[:len(session.prefixReqPath)-1], h.autoDirSlash)
			return
		}

		if data.CanCors {
			cors(w)
		}

		header(w, session.headers)

		if data.IsMutate {
			h.mutate(w, r, session, data)
			return
		}

		// archive
		if len(r.URL.RawQuery) >= 3 {
			switch r.URL.RawQuery[:3] {
			case "tar":
				h.tar(w, r, session, data)
				return
			case "tgz":
				h.tgz(w, r, session, data)
				return
			case "zip":
				h.zip(w, r, session, data)
				return
			}
		}
	}

	if h.applyMiddlewares(h.postMiddlewares, w, r, session, data) {
		return
	}

	// final process
	if session.wantJson {
		h.json(w, r, data)
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
	globalRestrictAccess = vhostCtx.restrictAccessDirs.mergePrefixMatched(globalRestrictAccess, util.HasFsPrefixDir, currentAlias.fs)
	globalRestrictAccess = util.InPlaceDedup(globalRestrictAccess)

	globalHeaders := p.GlobalHeaders
	globalHeaders = vhostCtx.headersUrls.mergePrefixMatched(globalHeaders, util.HasUrlPrefixDir, currentAlias.url)
	globalHeaders = vhostCtx.headersDirs.mergePrefixMatched(globalHeaders, util.HasFsPrefixDir, currentAlias.fs)

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

		globalAuth:    p.GlobalAuth || prefixMatched(p.AuthUrls, util.HasUrlPrefixDir, currentAlias.url) || prefixMatched(p.AuthDirs, util.HasFsPrefixDir, currentAlias.fs),
		authUrls:      filterSuccessor(p.AuthUrls, util.HasUrlPrefixDir, currentAlias.url),
		authUrlsUsers: vhostCtx.authUrlsUsers.filterSuccessor(true, util.HasUrlPrefixDir, currentAlias.url),
		authDirs:      filterSuccessor(p.AuthDirs, util.HasFsPrefixDir, currentAlias.fs),
		authDirsUsers: vhostCtx.authDirsUsers.filterSuccessor(true, util.HasFsPrefixDir, currentAlias.fs),

		globalUpload:    p.GlobalUpload || prefixMatched(p.UploadUrls, util.HasUrlPrefixDir, currentAlias.url) || prefixMatched(p.UploadDirs, util.HasFsPrefixDir, currentAlias.fs),
		uploadUrls:      filterSuccessor(p.UploadUrls, util.HasUrlPrefixDir, currentAlias.url),
		uploadUrlsUsers: vhostCtx.uploadUrlsUsers.filterSuccessor(true, util.HasUrlPrefixDir, currentAlias.url),
		uploadDirs:      filterSuccessor(p.UploadDirs, util.HasFsPrefixDir, currentAlias.fs),
		uploadDirsUsers: vhostCtx.uploadDirsUsers.filterSuccessor(true, util.HasFsPrefixDir, currentAlias.fs),

		globalMkdir:    p.GlobalMkdir || prefixMatched(p.MkdirUrls, util.HasUrlPrefixDir, currentAlias.url) || prefixMatched(p.MkdirDirs, util.HasFsPrefixDir, currentAlias.fs),
		mkdirUrls:      filterSuccessor(p.MkdirUrls, util.HasUrlPrefixDir, currentAlias.url),
		mkdirUrlsUsers: vhostCtx.mkdirUrlsUsers.filterSuccessor(true, util.HasUrlPrefixDir, currentAlias.url),
		mkdirDirs:      filterSuccessor(p.MkdirDirs, util.HasFsPrefixDir, currentAlias.fs),
		mkdirDirsUsers: vhostCtx.mkdirDirsUsers.filterSuccessor(true, util.HasFsPrefixDir, currentAlias.fs),

		globalDelete:    p.GlobalDelete || prefixMatched(p.DeleteUrls, util.HasUrlPrefixDir, currentAlias.url) || prefixMatched(p.DeleteDirs, util.HasFsPrefixDir, currentAlias.fs),
		deleteUrls:      filterSuccessor(p.DeleteUrls, util.HasUrlPrefixDir, currentAlias.url),
		deleteUrlsUsers: vhostCtx.deleteUrlsUsers.filterSuccessor(true, util.HasUrlPrefixDir, currentAlias.url),
		deleteDirs:      filterSuccessor(p.DeleteDirs, util.HasFsPrefixDir, currentAlias.fs),
		deleteDirsUsers: vhostCtx.deleteDirsUsers.filterSuccessor(true, util.HasFsPrefixDir, currentAlias.fs),

		globalArchive: p.GlobalArchive || prefixMatched(p.ArchiveUrls, util.HasUrlPrefixDir, currentAlias.url) || prefixMatched(p.ArchiveDirs, util.HasFsPrefixDir, currentAlias.fs),
		archiveUrls:   filterSuccessor(p.ArchiveUrls, util.HasUrlPrefixDir, currentAlias.url),
		archiveDirs:   filterSuccessor(p.ArchiveDirs, util.HasFsPrefixDir, currentAlias.fs),

		globalCors: p.GlobalCors || prefixMatched(p.CorsUrls, util.HasUrlPrefixDir, currentAlias.url) || prefixMatched(p.CorsDirs, util.HasFsPrefixDir, currentAlias.fs),
		corsUrls:   filterSuccessor(p.CorsUrls, util.HasUrlPrefixDir, currentAlias.url),
		corsDirs:   filterSuccessor(p.CorsDirs, util.HasFsPrefixDir, currentAlias.fs),

		globalRestrictAccess: globalRestrictAccess,
		restrictAccessUrls:   vhostCtx.restrictAccessUrls.filterSuccessor(false, util.HasUrlPrefixDir, currentAlias.url),
		restrictAccessDirs:   vhostCtx.restrictAccessDirs.filterSuccessor(false, util.HasFsPrefixDir, currentAlias.fs),

		globalHeaders: globalHeaders,
		headersUrls:   vhostCtx.headersUrls.filterSuccessor(false, util.HasUrlPrefixDir, currentAlias.url),
		headersDirs:   vhostCtx.headersDirs.filterSuccessor(false, util.HasFsPrefixDir, currentAlias.fs),

		shows:     vhostCtx.shows,
		showDirs:  vhostCtx.showDirs,
		showFiles: vhostCtx.showFiles,
		hides:     vhostCtx.hides,
		hideDirs:  vhostCtx.hideDirs,
		hideFiles: vhostCtx.hideFiles,

		vary: vhostCtx.vary,

		postMiddlewares: p.PostMiddlewares,
	}
	return h
}
