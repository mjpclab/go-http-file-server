package server

import (
	"../param"
	"../serverError"
	"../serverHandler"
	"../serverLog"
	"../tpl"
	"net/http"
	"text/template"
)

type Server struct {
	root     string
	key      string
	cert     string
	useTLS   bool
	listen   string
	tplFile  string
	tplObj   *template.Template
	aliases  map[string]string
	uploads  map[string]bool
	handlers map[string]http.Handler
	logger   *serverLog.Logger
}

var p *param.Param
var logger *serverLog.Logger

func init() {
	p = param.Parse()

	var err error
	logger, err = serverLog.NewLogger(p.AccessLog, p.ErrorLog)
	if !serverError.CheckFatal(err) {
		serverError.SetLogger(logger)
	}
}

func (s *Server) ListenAndServe() {
	var err error

	for urlPath, handler := range s.handlers {
		http.Handle(urlPath, handler)
		if len(urlPath) > 0 {
			http.Handle(urlPath+"/", handler)
		}
	}

	s.logger.LogAccess("Start to listen on " + s.listen)

	if s.useTLS {
		err = http.ListenAndServeTLS(s.listen, s.cert, s.key, nil)
	} else {
		err = http.ListenAndServe(s.listen, nil)
	}

	serverError.LogFatal(err)
}

func NewServer() *Server {
	useTLS := len(p.Key) > 0 && len(p.Cert) > 0

	listen := normalizePort(p.Listen, useTLS)

	tplObj, err := tpl.LoadPage(p.Template)
	serverError.LogError(err)

	aliases := p.Aliases
	uploads := p.Uploads
	handlers := map[string]http.Handler{}

	if _, hasAlias := aliases["/"]; !hasAlias {
		handlers["/"] = serverHandler.NewHandler(p.Root, "/", aliases, uploads, tplObj, logger)
	}

	for urlPath, fsPath := range p.Aliases {
		handlers[urlPath] = serverHandler.NewHandler(fsPath, urlPath, aliases, uploads, tplObj, logger)
	}

	return &Server{
		root:     p.Root,
		key:      p.Key,
		cert:     p.Cert,
		useTLS:   useTLS,
		listen:   listen,
		tplFile:  p.Template,
		tplObj:   tplObj,
		aliases:  aliases,
		uploads:  uploads,
		handlers: handlers,
		logger:   logger,
	}
}
