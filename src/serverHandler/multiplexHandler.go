package serverHandler

import (
	"mjpclab.dev/ghfs/src/param"
	"net/http"
)

type aliasWithHandler struct {
	alias   alias
	handler http.Handler
}

type multiplexHandler []aliasWithHandler

func (mux multiplexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, ah := range mux {
		if ah.alias.isMatch(r.URL.Path) || ah.alias.isPredecessorOf(r.URL.Path) {
			ah.handler.ServeHTTP(w, r)
			return
		}
	}

	defaultHandler.ServeHTTP(w, r)
}

func newMultiplexHandler(
	p *param.Param,
	vhostCtx *vhostContext,
) http.Handler {
	if len(p.Aliases) == 0 {
		return defaultHandler
	}

	aliases := newAliases(p.Aliases)

	if len(aliases) == 1 {
		alias, hasRootAlias := aliases.byUrlPath("/")
		if hasRootAlias {
			return newAliasHandler(p, vhostCtx, alias, aliases)
		}
	}

	aliasWithHandlers := make([]aliasWithHandler, len(aliases))
	for i, alias := range aliases {
		aliasWithHandlers[i] = aliasWithHandler{
			alias:   alias,
			handler: newAliasHandler(p, vhostCtx, alias, aliases),
		}
	}

	return multiplexHandler(aliasWithHandlers)
}
