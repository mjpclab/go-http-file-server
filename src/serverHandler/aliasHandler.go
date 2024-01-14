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

	globalAuth bool
	authUrls   []string
	authDirs   []string

	globalRestrictAccess []string
	restrictAccessUrls   pathStringsList
	restrictAccessDirs   pathStringsList

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

	inMiddlewares   []middleware.Middleware
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

	if h.applyMiddlewares(h.inMiddlewares, w, r, session, data) {
		return
	}

	if !session.allowAccess {
		if !h.applyMiddlewares(h.postMiddlewares, w, r, session, data) {
			h.page(w, r, session, data)
		}
		return
	}

	if session.needAuth {
		h.notifyAuth(w, r)
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

		globalAuth: p.GlobalAuth,
		authUrls:   p.AuthUrls,
		authDirs:   p.AuthDirs,

		globalRestrictAccess: globalRestrictAccess,
		restrictAccessUrls:   vhostCtx.restrictAccessUrls.filterSuccessor(util.HasUrlPrefixDir, currentAlias.url),
		restrictAccessDirs:   vhostCtx.restrictAccessDirs.filterSuccessor(util.HasFsPrefixDir, currentAlias.fs),

		globalHeaders: p.GlobalHeaders,
		headersUrls:   vhostCtx.headersUrls,
		headersDirs:   vhostCtx.headersDirs,

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

		shows:     vhostCtx.shows,
		showDirs:  vhostCtx.showDirs,
		showFiles: vhostCtx.showFiles,
		hides:     vhostCtx.hides,
		hideDirs:  vhostCtx.hideDirs,
		hideFiles: vhostCtx.hideFiles,

		vary: vhostCtx.vary,

		inMiddlewares:   p.InMiddlewares,
		postMiddlewares: p.PostMiddlewares,
	}
	return h
}
