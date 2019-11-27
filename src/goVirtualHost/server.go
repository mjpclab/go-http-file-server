package goVirtualHost

import (
	"crypto/tls"
	"net/http"
)

func newServer(useTLS bool) *server {
	server := &server{
		useTLS:       useTLS,
		vhosts:       vhosts{},
		defaultVhost: nil,

		httpServer: &http.Server{},
	}

	return server
}

func (server *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var vhost *vhost

	hostname := extractHostName(r.Host)

	for _, vh := range server.vhosts {
		if vh.matchHostName(hostname) {
			vhost = vh
			break
		}
	}

	if vhost == nil {
		vhost = server.defaultVhost
	}

	vhost.handler.ServeHTTP(w, r)
}

func (server *server) updateDefaultVhost() {
	for _, vh := range server.vhosts {
		if len(vh.hostNames) == 0 {
			server.defaultVhost = vh
			break
		}
	}

	if server.defaultVhost == nil {
		server.defaultVhost = server.vhosts[0]
	}
}

func (server *server) updateHttpServerTLSConfig() {
	var tlsConfig *tls.Config

	if server.useTLS {
		certs := []tls.Certificate{}

		for _, vhost := range server.vhosts {
			certs = append(certs, *vhost.cert)
		}

		tlsConfig = &tls.Config{
			Certificates: certs,
		}
		tlsConfig.BuildNameToCertificate()
	}

	server.httpServer.TLSConfig = tlsConfig
}

func (server *server) updateHttpServerHandler() {
	if len(server.vhosts) == 1 {
		server.httpServer.Handler = server.defaultVhost.handler
		return
	}

	server.httpServer.Handler = server
}

func (server *server) open(listener *listener) error {
	if server.useTLS {
		return server.httpServer.ServeTLS(listener.netListener, "", "")
	} else {
		return server.httpServer.Serve(listener.netListener)
	}
}

func (server *server) close() error {
	return server.httpServer.Close()
}
