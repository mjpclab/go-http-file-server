package goVirtualHost

import (
	"context"
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
	hostname := extractHostName(r.Host)
	for i := range server.vhosts {
		if server.vhosts[i].matchHostName(hostname) {
			server.vhosts[i].handler.ServeHTTP(w, r)
			return
		}
	}

	server.defaultVhost.handler.ServeHTTP(w, r)
}

func (server *server) getDefaultVhost() *vhost {
	for _, vh := range server.vhosts {
		if len(vh.hostNames) == 0 {
			return vh
		}
	}

	if len(server.vhosts) > 0 {
		return server.vhosts[0]
	}

	return nil
}

func (server *server) updateDefaultVhost() {
	server.defaultVhost = server.getDefaultVhost()
}

func (server *server) updateHttpServerTLSConfig() {
	var tlsConfig *tls.Config

	if server.useTLS {
		certs := []tls.Certificate{}

		for _, vhost := range server.vhosts {
			certs = append(certs, vhost.certs...)
		}

		tlsConfig = &tls.Config{
			Certificates: certs,
		}
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

func (server *server) shutdown(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}

func (server *server) close() error {
	return server.httpServer.Close()
}
