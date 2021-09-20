package vhostHandler

import (
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"../user"
	"net/http"
)

type VhostHandler struct {
	p            *param.Param
	logger       *serverLog.Logger
	errorHandler *serverErrHandler.ErrHandler
	theme        tpl.Theme
	Handler      http.Handler
}

func NewHandler(
	p *param.Param,
	logger *serverLog.Logger,
	errorHandler *serverErrHandler.ErrHandler,
	theme tpl.Theme,
) *VhostHandler {
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

	// register handlers
	handlers := map[string]http.Handler{}
	for urlPath, fsPath := range p.Aliases {
		emptyHandlerRoot := p.EmptyRoot && urlPath == "/"
		handlers[urlPath] = serverHandler.NewHandler(fsPath, emptyHandlerRoot, urlPath, p, users, theme, logger, errorHandler)
	}

	var handler http.Handler
	if len(handlers) == 1 {
		handler = handlers["/"]
	}
	if handler == nil {
		serveMux := http.NewServeMux()
		for urlPath, urlHandler := range handlers {
			serveMux.Handle(urlPath, urlHandler)
			if len(urlPath) > 1 {
				serveMux.Handle(urlPath+"/", urlHandler)
			}
		}
		handler = serveMux
	}

	vhostHandler := &VhostHandler{
		p:            p,
		logger:       logger,
		errorHandler: errorHandler,
		Handler:      handler,
	}

	return vhostHandler
}

func (m *VhostHandler) ReOpenLog() {
	errors := m.logger.ReOpen()
	serverErrHandler.CheckError(errors...)
}

func (m *VhostHandler) Close() {
	m.logger.Close()
}
