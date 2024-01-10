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

type pathStrings struct {
	path    string
	strings []string
}

type aliasHandler struct {
	root         string
	emptyRoot    bool
	autoDirSlash int
	hsts         bool
	hstsMaxAge   string
	toHttps      bool
	toHttpsPort  string // with prefix ":"
	defaultSort  string
	aliasPrefix  string

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
			h.page(w, r, data)
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
		h.page(w, r, data)
	}
}

func newAliasHandler(
	p *param.Param,
	vhostCtx *vhostContext,
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

	h := &aliasHandler{
		root:         currentAlias.fs,
		emptyRoot:    emptyRoot,
		autoDirSlash: p.AutoDirSlash,
		hsts:         p.Hsts,
		hstsMaxAge:   strconv.Itoa(p.HstsMaxAge),
		toHttps:      p.ToHttps,
		toHttpsPort:  p.ToHttpsPort,
		defaultSort:  p.DefaultSort,
		aliasPrefix:  currentAlias.url,

		users:  vhostCtx.users,
		theme:  vhostCtx.theme,
		logger: vhostCtx.logger,

		dirIndexes: p.DirIndexes,
		aliases:    aliases,

		globalAuth: p.GlobalAuth,
		authUrls:   p.AuthUrls,
		authDirs:   p.AuthDirs,

		restrictAccess:       vhostCtx.restrictAccess,
		globalRestrictAccess: p.GlobalRestrictAccess,
		restrictAccessUrls:   vhostCtx.restrictAccessUrls,
		restrictAccessDirs:   vhostCtx.restrictAccessDirs,

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
