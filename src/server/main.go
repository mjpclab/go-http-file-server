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
	root    string
	key     string
	cert    string
	useTLS  bool
	listen  string
	tplFile string
	tplObj  *template.Template
	handler http.Handler
}

func (s *Server) ListenAndServe() {
	var err error

	if s.useTLS {
		err = http.ListenAndServeTLS(s.listen, s.cert, s.key, s.handler)
	} else {
		err = http.ListenAndServe(s.listen, s.handler)
	}

	serverError.CheckFatal(err)
}

func NewServer() *Server {
	p := param.Parse()

	useTLS := len(p.Key) > 0 && len(p.Cert) > 0

	listen := normalizePort(p.Listen, useTLS)

	tplObj := tpl.LoadPage(p.Template)

	handler := serverHandler.NewHandler(p.Root, tplObj)

	return &Server{
		root:    p.Root,
		key:     p.Key,
		cert:    p.Cert,
		useTLS:  useTLS,
		listen:  listen,
		tplFile: p.Template,
		tplObj:  tplObj,
		handler: handler,
	}
}
