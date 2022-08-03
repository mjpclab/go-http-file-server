package serverHandler

import (
	"../param"
	"net/http"
)

type aliasWithHandler struct {
	alias   alias
	handler http.Handler
}

type multiplexHandler struct {
	aliasWithHandlers []aliasWithHandler
}

func (mux multiplexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, ah := range mux.aliasWithHandlers {
		if ah.alias.isMatch(r.URL.Path) || ah.alias.isPredecessorOf(r.URL.Path) {
			ah.handler.ServeHTTP(w, r)
			return
		}
	}

	defaultHandler.ServeHTTP(w, r)
}

func newMultiplexHandler(
	p *param.Param,
	ap *aliasParam,
) http.Handler {
	if len(p.Aliases) == 0 {
		return defaultHandler
	}

	aliases := newAliases(p.Aliases)

	if len(aliases) == 1 {
		alias, hasRootAlias := aliases.byUrlPath("/")
		if hasRootAlias {
			return newAliasHandler(p, ap, alias, aliases)
		}
	}

	aliasWithHandlers := make([]aliasWithHandler, len(aliases))
	for i, alias := range aliases {
		aliasWithHandlers[i] = aliasWithHandler{
			alias:   alias,
			handler: newAliasHandler(p, ap, alias, aliases),
		}
	}
	return multiplexHandler{aliasWithHandlers}
}