package serverHandler

import (
	"../param"
	"../serverLog"
	"../tpl"
	"../user"
	"../util"
	"net/http"
	"strings"
)

type aliasHandler struct {
	alias   alias
	handler http.Handler
}

type multiplexer struct {
	aliasHandlers []aliasHandler
}

func (mux multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawReqPath := util.CleanUrlPath(r.URL.Path)
	for _, aliasHandler := range mux.aliasHandlers {
		if aliasHandler.alias.isMatch(rawReqPath) || aliasHandler.alias.isPredecessorOf(rawReqPath) {
			aliasHandler.handler.ServeHTTP(w, r)
			return
		}
	}

	defaultHandler.ServeHTTP(w, r)
}

func NewMultiplexer(
	p *param.Param,
	users user.List,
	theme tpl.Theme,
	logger *serverLog.Logger,
) http.Handler {
	if len(p.Aliases) == 0 {
		return defaultHandler
	}

	aliases := newAliases(p.Aliases)
	restrictAccessUrls := newRestrictAccesses(p.RestrictAccessUrls)
	restrictAccessDirs := newRestrictAccesses(p.RestrictAccessDirs)
	headersUrls := newPathHeaders(p.HeadersUrls)
	headersDirs := newPathHeaders(p.HeadersDirs)

	restrictAccess := hasRestrictAccess(p.GlobalRestrictAccess, restrictAccessUrls, restrictAccessDirs)
	pageVaryV1 := "Accept-Encoding"
	contentVaryV1 := ""
	if restrictAccess {
		pageVaryV1 += ", Referer"
		contentVaryV1 = "Referer"
	}
	pageVary := strings.ToLower(pageVaryV1)
	contentVary := strings.ToLower(contentVaryV1)

	if len(aliases) == 1 {
		alias, hasRootAlias := aliases.byUrlPath("/")
		if hasRootAlias {
			return newHandler(
				p, alias.fs, alias.url, aliases,
				restrictAccessUrls, restrictAccessDirs,
				headersUrls, headersDirs,
				users, theme, logger,
				restrictAccess,
				pageVaryV1, pageVary, contentVaryV1, contentVary,
			)
		}
	}

	aliasHandlers := make([]aliasHandler, len(aliases))
	for i, alias := range aliases {
		aliasHandlers[i] = aliasHandler{
			alias: alias,
			handler: newHandler(
				p, alias.fs, alias.url, aliases,
				restrictAccessUrls, restrictAccessDirs,
				headersUrls, headersDirs,
				users, theme, logger,
				restrictAccess,
				pageVaryV1, pageVary, contentVaryV1, contentVary,
			),
		}
	}
	return multiplexer{aliasHandlers}
}
