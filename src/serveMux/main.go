package serveMux

import (
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"net/http"
)

type ServeMux struct {
	http.ServeMux
	logger     *serverLog.Logger
	errHandler *serverErrHandler.ErrHandler
}

func NewServeMux(
	p *param.Param,
	logger *serverLog.Logger,
	errorHandler *serverErrHandler.ErrHandler,
) *ServeMux {
	tplObj, err := tpl.LoadPage(p.Template)
	errorHandler.LogError(err)

	aliases := p.Aliases
	handlers := map[string]http.Handler{}

	if _, hasRootAlias := aliases["/"]; !hasRootAlias {
		handlers["/"] = serverHandler.NewHandler(p.Root, "/", p, tplObj, logger, errorHandler)
	}

	for urlPath, fsPath := range p.Aliases {
		handlers[urlPath] = serverHandler.NewHandler(fsPath, urlPath, p, tplObj, logger, errorHandler)
	}

	// create ServeMux
	serveMux := &ServeMux{
		logger:     logger,
		errHandler: errorHandler,
	}
	for urlPath, handler := range handlers {
		serveMux.Handle(urlPath, handler)
		if len(urlPath) > 1 {
			serveMux.Handle(urlPath+"/", handler)
		}
	}

	return serveMux
}
