package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/serverLog"
	"net/http"
)

type preprocessHandler struct {
	logger         *serverLog.Logger
	preMiddlewares []middleware.Middleware
	nextHandler    http.Handler
}

func (pph preprocessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(pph.logger, r)

	if len(pph.preMiddlewares) > 0 {
		middlewareContext := &middleware.Context{
			PrefixReqPath: r.URL.RawPath, // init by pathTransformHandler
			VhostReqPath:  r.URL.Path,
			Logger:        pph.logger,
		}
		for i := range pph.preMiddlewares {
			processResult := pph.preMiddlewares[i](w, r, middlewareContext)
			if processResult == middleware.Processed {
				return
			} else if processResult == middleware.SkipRests {
				break
			}
		}
	}

	pph.nextHandler.ServeHTTP(w, r)
}

func newPreprocessHandler(logger *serverLog.Logger, preMiddlewares []middleware.Middleware, nextHandler http.Handler) http.Handler {
	return preprocessHandler{
		logger:         logger,
		preMiddlewares: preMiddlewares,
		nextHandler:    nextHandler,
	}
}
