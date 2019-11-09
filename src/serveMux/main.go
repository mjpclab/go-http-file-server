package serveMux

import (
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"../user"
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
	users := user.NewUsers()
	for username, password := range p.Users {
		users.Add(username, password)
	}

	tplObj, err := tpl.LoadPage(p.Template)
	errorHandler.LogError(err)

	aliases := p.Aliases
	handlers := map[string]http.Handler{}

	if _, hasRootAlias := aliases["/"]; !hasRootAlias {
		handlers["/"] = serverHandler.NewHandler(p.Root, "/", p, users, tplObj, logger, errorHandler)
	}

	for urlPath, fsPath := range p.Aliases {
		handlers[urlPath] = serverHandler.NewHandler(fsPath, urlPath, p, users, tplObj, logger, errorHandler)
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
