package goVirtualHost

import (
	"context"
	"crypto/tls"
	"net/http"
)

func newServeable(useTLS bool) *serveable {
	return &serveable{
		useTLS:       useTLS,
		vhosts:       vhosts{},
		defaultVhost: nil,

		server: &http.Server{},
	}
}

func (serveable *serveable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostname := extractHostName(r.Host)
	for i := range serveable.vhosts {
		if serveable.vhosts[i].matchHostName(hostname) {
			serveable.vhosts[i].handler.ServeHTTP(w, r)
			return
		}
	}

	serveable.defaultVhost.handler.ServeHTTP(w, r)
}

func (serveable *serveable) getDefaultVhost() *vhost {
	for _, vh := range serveable.vhosts {
		if len(vh.hostNames) == 0 {
			return vh
		}
	}

	if len(serveable.vhosts) > 0 {
		return serveable.vhosts[0]
	}

	return nil
}

func (serveable *serveable) updateDefaultVhost() {
	serveable.defaultVhost = serveable.getDefaultVhost()
}

func (serveable *serveable) updateHttpServerTLSConfig() {
	if !serveable.useTLS {
		return
	}

	certs := make([]tls.Certificate, 0, len(serveable.vhosts))

	for _, vhost := range serveable.vhosts {
		certs = append(certs, vhost.certs...)
	}

	serveable.server.TLSConfig = &tls.Config{
		Certificates: certs,
	}
}

func (serveable *serveable) updateHttpServerHandler() {
	if len(serveable.vhosts) == 1 {
		serveable.server.Handler = serveable.defaultVhost.handler
		return
	}

	serveable.server.Handler = serveable
}

func (serveable *serveable) open(l *listenable) error {
	if serveable.useTLS {
		return serveable.server.ServeTLS(l.listener, "", "")
	} else {
		return serveable.server.Serve(l.listener)
	}
}

func (serveable *serveable) shutdown(ctx context.Context) error {
	return serveable.server.Shutdown(ctx)
}

func (serveable *serveable) close() error {
	return serveable.server.Close()
}
