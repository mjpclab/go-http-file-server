package serverHandler

import (
	"../param"
	"../serverErrHandler"
	"../serverLog"
	"../tpl"
	"../user"
	"../util"
	"net/http"
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
	errHandler *serverErrHandler.ErrHandler,
) http.Handler {
	aliases := newAliases(p.Aliases, p.Binds)

	if len(aliases) == 0 {
		return defaultHandler
	}

	if len(aliases) == 1 {
		alias, hasRootAlias := aliases.byUrlPath("/")
		if hasRootAlias {
			return newHandler(p, alias.fsPath(), alias.urlPath(), aliases, users, theme, logger, errHandler)
		}
	}

	aliasHandlers := make([]aliasHandler, len(aliases))
	for i, alias := range aliases {
		aliasHandlers[i] = aliasHandler{
			alias:   alias,
			handler: newHandler(p, alias.fsPath(), alias.urlPath(), aliases, users, theme, logger, errHandler),
		}
	}
	return multiplexer{aliasHandlers}
}
