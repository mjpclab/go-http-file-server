package server

import (
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"net/http"
)

type Server struct {
	key        string
	cert       string
	useTLS     bool
	listen     string
	handlers   map[string]http.Handler
	logger     *serverLog.Logger
	errHandler *serverErrHandler.ErrHandler
}

func (s *Server) ListenAndServe() {
	var err error

	for urlPath, handler := range s.handlers {
		http.Handle(urlPath, handler)
		if len(urlPath) > 0 {
			http.Handle(urlPath+"/", handler)
		}
	}

	s.logger.LogAccessString("start to listen on " + s.listen)

	if s.useTLS {
		err = http.ListenAndServeTLS(s.listen, s.cert, s.key, nil)
	} else {
		err = http.ListenAndServe(s.listen, nil)
	}

	s.errHandler.LogError(err)
}

func NewServer(p *param.Param) *Server {
	logger, err := serverLog.NewLogger(p.AccessLog, p.ErrorLog)
	serverErrHandler.CheckFatal(err)

	errorHandler := serverErrHandler.NewErrHandler(logger)

	useTLS := len(p.Key) > 0 && len(p.Cert) > 0

	listen := normalizePort(p.Listen, useTLS)

	tplObj, err := tpl.LoadPage(p.Template)
	errorHandler.LogError(err)

	aliases := p.Aliases
	handlers := map[string]http.Handler{}

	if _, hasAlias := aliases["/"]; !hasAlias {
		handlers["/"] = serverHandler.NewHandler(p.Root, "/", p, tplObj, logger, errorHandler)
	}

	for urlPath, fsPath := range p.Aliases {
		handlers[urlPath] = serverHandler.NewHandler(fsPath, urlPath, p, tplObj, logger, errorHandler)
	}

	return &Server{
		key:        p.Key,
		cert:       p.Cert,
		useTLS:     useTLS,
		listen:     listen,
		handlers:   handlers,
		logger:     logger,
		errHandler: errorHandler,
	}
}

func (s *Server) Close() {
	s.logger.Close()
}
