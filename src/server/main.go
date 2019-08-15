package server

import (
	"../param"
	"../serverError"
	"../serverHandler"
	"../tpl"
	"net/http"
	"path"
	"text/template"
)

type Server struct {
	root    string
	listen  string
	key     string
	cert    string
	tplFile string
	tplObj  *template.Template
	handler http.Handler
}

func (s *Server) ListenAndServe() {
	var err error

	if len(s.key) == 0 || len(s.cert) == 0 {
		err = http.ListenAndServe(s.listen, s.handler)
	} else {
		err = http.ListenAndServeTLS(s.listen, s.cert, s.key, s.handler)
	}

	serverError.CheckFatal(err)
}

func NewServer() *Server {
	p := param.Parse()

	var tplObj *template.Template
	var err error
	if len(p.Template) > 0 {
		tplObj, err = template.New(path.Base(p.Template)).ParseFiles(p.Template)
		serverError.CheckFatal(err)
	}
	if err != nil || len(p.Template) == 0 {
		tplObj = tpl.Page
	}

	handler := serverHandler.NewHandler(p.Root, tplObj)

	return &Server{
		root:    p.Root,
		listen:  p.Listen,
		key:     p.Key,
		cert:    p.Cert,
		tplFile: p.Template,
		tplObj:  tplObj,
		handler: handler,
	}
}
