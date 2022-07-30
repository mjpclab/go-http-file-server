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

type aliasWithHandler struct {
	alias   alias
	handler http.Handler
}

type multiplexHandler struct {
	aliasWithHandlers []aliasWithHandler
}

func (mux multiplexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawReqPath := util.CleanUrlPath(r.URL.Path)
	for _, aAndH := range mux.aliasWithHandlers {
		if aAndH.alias.isMatch(rawReqPath) || aAndH.alias.isPredecessorOf(rawReqPath) {
			aAndH.handler.ServeHTTP(w, r)
			return
		}
	}

	defaultHandler.ServeHTTP(w, r)
}

func NewMultiplexHandler(
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
			return newAliasHandler(
				p, alias.fs, alias.url, aliases,
				restrictAccessUrls, restrictAccessDirs,
				headersUrls, headersDirs,
				users, theme, logger,
				restrictAccess,
				pageVaryV1, pageVary, contentVaryV1, contentVary,
			)
		}
	}

	aliasHandlers := make([]aliasWithHandler, len(aliases))
	for i, alias := range aliases {
		aliasHandlers[i] = aliasWithHandler{
			alias: alias,
			handler: newAliasHandler(
				p, alias.fs, alias.url, aliases,
				restrictAccessUrls, restrictAccessDirs,
				headersUrls, headersDirs,
				users, theme, logger,
				restrictAccess,
				pageVaryV1, pageVary, contentVaryV1, contentVary,
			),
		}
	}
	return multiplexHandler{aliasHandlers}
}
