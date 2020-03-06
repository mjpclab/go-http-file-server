package vhostMux

import (
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"../user"
	"net/http"
)

type VhostMux struct {
	p            *param.Param
	logger       *serverLog.Logger
	errorHandler *serverErrHandler.ErrHandler
	ServeMux     *http.ServeMux
}

func NewServeMux(
	p *param.Param,
	logger *serverLog.Logger,
	errorHandler *serverErrHandler.ErrHandler,
) *VhostMux {
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

	// template
	pageTpl, err := tpl.LoadPage(p.Template)
	errorHandler.LogError(err)

	// register handlers
	aliases := p.Aliases
	_, hasRootAlias := aliases["/"]
	emptyRoot := false
	if !hasRootAlias {
		aliases["/"] = p.Root
		emptyRoot = p.EmptyRoot
	}

	handlers := map[string]http.Handler{}
	for urlPath, fsPath := range aliases {
		emptyHandlerRoot := emptyRoot && urlPath == "/"
		handlers[urlPath] = serverHandler.NewHandler(fsPath, emptyHandlerRoot, urlPath, p, users, pageTpl, logger, errorHandler)
	}

	// create ServeMux
	serveMux := &http.ServeMux{}
	for urlPath, handler := range handlers {
		serveMux.Handle(urlPath, handler)
		if len(urlPath) > 1 {
			serveMux.Handle(urlPath+"/", handler)
		}
	}

	vhostMux := &VhostMux{
		p:            p,
		logger:       logger,
		errorHandler: errorHandler,
		ServeMux:     serveMux,
	}

	return vhostMux
}

func (m *VhostMux) ReOpenLog() {
	errors := m.logger.ReOpen()
	serverErrHandler.CheckError(errors...)
}

func (m *VhostMux) Close() {
	m.logger.Close()
}
