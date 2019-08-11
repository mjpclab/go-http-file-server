package server

import (
	"../handler"
	"../param"
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	root    string
	listen  string
	key     string
	cert    string
	handler http.Handler
}

func (s *Server) ListenAndServe() {
	var err error

	if len(s.key) == 0 || len(s.cert) == 0 {
		err = http.ListenAndServe(s.listen, s.handler)
	} else {
		err = http.ListenAndServeTLS(s.listen, s.cert, s.key, s.handler)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewServer() *Server {
	p := param.Parse()
	h := handler.NewHandler(p.Root)
	return &Server{
		root:    p.Root,
		listen:  p.Listen,
		key:     p.Key,
		cert:    p.Cert,
		handler: h,
	}
}
