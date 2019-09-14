package server

import (
	"../param"
	"../serverErrHandler"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"net"
	"net/http"
	"os"
)

type Server struct {
	key         string
	cert        string
	useTLS      bool
	listenProto string
	listenAddr  string
	handlers    map[string]http.Handler
	logger      *serverLog.Logger
	listener    net.Listener
	errHandler  *serverErrHandler.ErrHandler
}

func (s *Server) openTransListener() (err error) {
	s.listener, err = net.Listen(s.listenProto, s.listenAddr)
	s.errHandler.LogError(err)

	return
}

func (s *Server) closeTransListener() (err error) {
	if s.listener == nil {
		return
	}

	err = s.listener.Close()
	s.listener = nil
	s.errHandler.LogError(err)

	return
}

func (s *Server) ListenAndServe() {
	var err error

	for urlPath, handler := range s.handlers {
		http.Handle(urlPath, handler)
		if len(urlPath) > 0 {
			http.Handle(urlPath+"/", handler)
		}
	}

	s.logger.LogAccessString("start to listen on " + s.listenProto + ": " + s.listenAddr)

	if s.openTransListener() != nil {
		return
	}

	server := &http.Server{}
	if s.useTLS {
		err = server.ServeTLS(s.listener, s.cert, s.key)
	} else {
		err = server.Serve(s.listener)
	}
	s.errHandler.LogError(err)

	s.closeTransListener()
}

func NewServer(p *param.Param) *Server {
	logger, err := serverLog.NewLogger(p.AccessLog, p.ErrorLog)
	serverErrHandler.CheckFatal(err)

	errorHandler := serverErrHandler.NewErrHandler(logger)

	useTLS := len(p.Key) > 0 && len(p.Cert) > 0

	listenProto, listenAddr := splitListen(p.Listen, useTLS)
	if listenProto == "unix" {
		sockFile, _ := os.Lstat(listenAddr)
		if sockFile != nil && (sockFile.Mode()&os.ModeSocket != 0) {
			os.Remove(listenAddr)
		}
	}

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
		key:         p.Key,
		cert:        p.Cert,
		useTLS:      useTLS,
		listenProto: listenProto,
		listenAddr:  listenAddr,
		handlers:    handlers,
		logger:      logger,
		errHandler:  errorHandler,
	}
}

func (s *Server) Close() {
	s.logger.Close()
	s.closeTransListener()
}
