package goVirtualHost

import (
	"net/http"
	"strings"
)

func newVhost(certs certs, hostNames []string, handler http.Handler) *vhost {
	vhost := &vhost{
		certs:     certs,
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
		if len(hostname) > 1 {
			if hostname[0] == '.' && strings.HasSuffix(reqHostName, hostname) {
				return true
			} else if hostname[len(hostname)-1] == '.' && strings.HasPrefix(reqHostName, hostname) {
				return true
			}
		}
	}
	return false
}
