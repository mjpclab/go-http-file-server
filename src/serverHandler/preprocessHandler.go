package serverHandler

import (
	"mjpclab.dev/ghfs/src/middleware"
	"mjpclab.dev/ghfs/src/serverCompress"
	"mjpclab.dev/ghfs/src/serverLog"
	"net/http"
	"strings"
)

type preprocessHandler struct {
	logger         *serverLog.Logger
	preMiddlewares []middleware.Middleware
	nextHandler    http.Handler
}

func (pph preprocessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(pph.logger, r)

	rw := serverCompress.NewResponseWriter(w, r)

	if len(pph.preMiddlewares) > 0 {
		prefixReqPath := r.RequestURI // init by pathTransformHandler
		if qsIndex := strings.IndexByte(prefixReqPath, '?'); qsIndex >= 0 {
			prefixReqPath = prefixReqPath[:qsIndex]
		}
		middlewareContext := &middleware.Context{
			PrefixReqPath: prefixReqPath,
			VhostReqPath:  r.URL.Path,
			Logger:        pph.logger,
		}
		for i := range pph.preMiddlewares {
			processResult := pph.preMiddlewares[i](rw, r, middlewareContext)
			if processResult == middleware.Outputted {
				rw.Close()
				return
			} else if processResult == middleware.SkipRests {
				break
			}
		}
	}

	pph.nextHandler.ServeHTTP(rw, r)
	rw.Close()
}

func newPreprocessHandler(logger *serverLog.Logger, preMiddlewares []middleware.Middleware, nextHandler http.Handler) http.Handler {
	return preprocessHandler{
		logger:         logger,
		preMiddlewares: preMiddlewares,
		nextHandler:    nextHandler,
	}
}
