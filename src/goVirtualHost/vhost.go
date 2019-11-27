package goVirtualHost

import (
	"crypto/tls"
	"net/http"
	"strings"
)

func newVhost(cert *tls.Certificate, hostNames []string, handler http.Handler) *vhost {
	vhost := &vhost{
		cert:      cert,
		hostNames: hostNames,
		handler:   handler,
	}

	return vhost
}

func (v *vhost) matchHostName(name string) bool {
	reqHostName := strings.ToLower(name)
	for _, hostname := range v.hostNames {
		if hostname == reqHostName {
			return true
		}
		if len(hostname) > 0 && hostname[0] == '.' && strings.HasSuffix(reqHostName, hostname) {
			return true
		}
	}
	return false
}
