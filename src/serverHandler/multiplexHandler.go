package serverHandler

import (
	"mjpclab.dev/ghfs/src/param"
	"net/http"
)

type multiplexHandler []*aliasHandler

func (mux multiplexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range mux {
		if h.isMatch(r.URL.Path) || h.isPredecessorOf(r.URL.Path) {
			h.ServeHTTP(w, r)
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
		if alias, hasRootAlias := aliases.byUrlPath("/"); hasRootAlias {
			return newAliasHandler(p, vhostCtx, alias, aliases)
		}
	}

	mux := make(multiplexHandler, len(aliases))
	for i, alias := range aliases {
		mux[i] = newAliasHandler(p, vhostCtx, alias, aliases)
	}
	return mux
}
