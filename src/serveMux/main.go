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
	for _, u := range p.UsersPlain {
		errorHandler.LogError(users.AddPlain(u.Username, u.Password))
	}
	for _, u := range p.UsersBase64 {
		errorHandler.LogError(users.AddBase64(u.Username, u.Password))
	}
	for _, u := range p.UsersMd5 {
		errorHandler.LogError(users.AddMd5(u.Username, u.Password))
	}
	for _, u := range p.UsersSha1 {
		errorHandler.LogError(users.AddSha1(u.Username, u.Password))
	}
	for _, u := range p.UsersSha256 {
		errorHandler.LogError(users.AddSha256(u.Username, u.Password))
	}
	for _, u := range p.UsersSha512 {
		errorHandler.LogError(users.AddSha512(u.Username, u.Password))
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
