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
		server:       &http.Server{},
	}
}

func (serveable *serveable) lookupVhost(hostname string) *vhost {
	if len(serveable.vhosts) == 1 {
		return serveable.vhosts[0]
	}

	for _, vh := range serveable.vhosts {
		if vh.matchHostName(hostname) {
			return vh
		}
	}

	return serveable.defaultVhost
}

func (serveable *serveable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostname := extractHostName(r.Host)
	vh := serveable.lookupVhost(hostname)
	vh.handler.ServeHTTP(w, r)
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

func (serveable *serveable) updateHttpServerHandler() {
	if len(serveable.vhosts) == 1 {
		serveable.server.Handler = serveable.defaultVhost.handler
		return
	}

	serveable.server.Handler = serveable
}

func (serveable *serveable) loadCertificates() (errs []error) {
	if !serveable.useTLS {
		return
	}

	for _, vh := range serveable.vhosts {
		es := vh.loadCertificates()
		errs = append(errs, es...)
	}

	return
}

func (serveable *serveable) updateHttpServerTLSConfig() (errs []error) {
	if !serveable.useTLS {
		return
	}

	serveable.server.TLSConfig = &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			vh := serveable.lookupVhost(hello.ServerName)
			return vh.lookupCertificate(hello)
		},
	}

	return
}

func (serveable *serveable) init() (errs []error) {
	serveable.updateDefaultVhost()
	serveable.updateHttpServerHandler()
	errs = serveable.loadCertificates()
	serveable.updateHttpServerTLSConfig()
	return
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
