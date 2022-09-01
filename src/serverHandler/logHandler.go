package serverHandler

import (
	"mjpclab.dev/ghfs/src/serverLog"
	"net/http"
)

type logHandler struct {
	logger      *serverLog.Logger
	nextHandler http.Handler
}

func (l logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(l.logger, r)
	l.nextHandler.ServeHTTP(w, r)
}

func newLogHandler(logger *serverLog.Logger, nextHandler http.Handler) http.Handler {
	return logHandler{
		logger:      logger,
		nextHandler: nextHandler,
	}
}
