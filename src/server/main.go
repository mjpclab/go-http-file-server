package server

import (
	"../param"
	"../serverError"
	"../serverHandler"
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
}

func (s *Server) ListenAndServe() {
	var err error

	for urlPath, handler := range s.handlers {
		http.Handle(urlPath, handler)
		if len(urlPath) > 0 {
			http.Handle(urlPath+"/", handler)
		}
	}

	if s.useTLS {
		err = http.ListenAndServeTLS(s.listen, s.cert, s.key, nil)
	} else {
		err = http.ListenAndServe(s.listen, nil)
	}

	serverError.CheckFatal(err)
}

func NewServer() *Server {
	p := param.Parse()

	useTLS := len(p.Key) > 0 && len(p.Cert) > 0

	listen := normalizePort(p.Listen, useTLS)

	tplObj := tpl.LoadPage(p.Template)

	aliases := p.Aliases
	uploads := p.Uploads
	handlers := map[string]http.Handler{}

	if _, hasAlias := aliases["/"]; !hasAlias {
		handlers["/"] = serverHandler.NewHandler(p.Root, "/", aliases, uploads, tplObj)
	}

	for urlPath, fsPath := range p.Aliases {
		handlers[urlPath] = serverHandler.NewHandler(fsPath, urlPath, aliases,uploads, tplObj)
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
	}
}
