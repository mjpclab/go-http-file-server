package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/param"
	"net/http"
)

type aliasWithHandler struct {
	alias   alias
	handler http.Handler
}

type multiplexHandler struct {
	preMiddlewares    []middleware.Middleware
	aliasWithHandlers []aliasWithHandler
}

func (mux multiplexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(mux.preMiddlewares) > 0 {
		middlewareContext := &middleware.Context{
			PrefixReqPath: r.URL.RawPath, // init by pathTransformHandler
			VhostReqPath:  r.URL.Path,
		}
		for i := range mux.preMiddlewares {
			processResult := mux.preMiddlewares[i](w, r, middlewareContext)
			if processResult == middleware.Processed {
				return
			} else if processResult == middleware.SkipRests {
				break
			}
		}
	}

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

	if len(aliases) == 1 && len(p.PreMiddlewares) == 0 {
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

	return multiplexHandler{
		preMiddlewares:    p.PreMiddlewares,
		aliasWithHandlers: aliasWithHandlers,
	}
}
